package dusk

import (
	"github.com/go-gl/mathgl/mgl32"

	"github.com/WhoBrokeTheBuild/GoDusk/m32"
)

// ForceMode is the type of force to apply to a RigidBody
type ForceMode int

const (
	// ConstantForce adds a continuous force, using mass
	ConstantForce ForceMode = iota
	// Acceleration adds a continuous force, ignoring mass
	Acceleration = iota
	// Impulse adds an instant force, using mass
	Impulse = iota
	// VelocityChange adds an instant force, ignoring mass
	VelocityChange = iota
)

type ICollider interface {
}

type SphereCollider struct {
	Radius float32
}

type BoxCollider struct {
	Origin mgl32.Vec3
	Size   mgl32.Vec3
}

// RigidBody is a physics body implemented with Rigid Body dynamics
type RigidBody struct {
	Parent          IActor
	Collider        ICollider
	Mass            float32
	Elasticity      float32
	Restitution     float32
	Position        mgl32.Vec3
	AngularVelocity mgl32.Vec3
	Velocity        mgl32.Vec3
	Acceleration    mgl32.Vec3
}

// NewRigidBody creates a new RigidBody with appropriate defaults
func NewRigidBody(parent IActor) *RigidBody {
	return &RigidBody{
		Parent:       parent,
		Mass:         1.0,
		Elasticity:   0.8,
		Restitution:  0.05,
		Velocity:     mgl32.Vec3{0, 0, 0},
		Acceleration: mgl32.Vec3{0, 0, 0},
	}
}

// Delete frees up resources
func (rb *RigidBody) Delete() {
	rb.Parent = nil
}

// ApplyForce adds a force to the object, how it is added depenends on the mode
func (rb *RigidBody) ApplyForce(force mgl32.Vec3, mode ForceMode) {
	switch mode {
	case ConstantForce:
		force = force.Mul(1.0 / rb.Mass)
		rb.Acceleration = rb.Acceleration.Add(force)
	case Acceleration:
		rb.Acceleration = rb.Acceleration.Add(force)
	case Impulse:
		force = force.Mul(1.0 / rb.Mass)
		rb.Velocity = rb.Velocity.Add(force)
	case VelocityChange:
		rb.Velocity = rb.Velocity.Add(force)
	}
}

// Update updates a rigidbody
func (rb *RigidBody) Update(ctx *UpdateContext) {
	rb.Parent.Transform().Position = rb.Position
	rb.Position = rb.Parent.Transform().Position.Add(rb.Velocity.Mul(ctx.DeltaTime))
	rb.Velocity = rb.Velocity.Add(rb.Acceleration.Mul(float32(ctx.DeltaTime)))

	// Min Velocity
	if rb.Velocity.Len() < rb.Restitution {
		rb.Velocity = mgl32.Vec3{0, 0, 0}
	}
}

// CheckCollide checks whether a collision has occurred between the two rigid bodies
func (rb *RigidBody) CheckCollide(other *RigidBody) {
	pos := rb.Parent.Transform().Position
	otherPos := other.Parent.Transform().Position

	checkSphereSphere := func(posA, posB mgl32.Vec3, colA, colB SphereCollider) bool {
		dist := DistanceSquared(posA, posB)
		if dist < (colA.Radius+colB.Radius)*(colA.Radius+colB.Radius) {
			return true
		}
		return false
	}

	checkBoxBox := func(posA, posB mgl32.Vec3, colA, colB BoxCollider) bool {
		aMin := posA.Sub(colA.Origin)
		aMax := aMin.Add(colA.Size)
		bMin := posB.Sub(colB.Origin)
		bMax := bMin.Add(colB.Size)

		if aMin.X() > bMax.X() || bMin.X() > aMax.X() {
			return false
		}

		if aMin.Y() > bMax.Y() || bMin.Y() > aMax.Y() {
			return false
		}

		if aMin.Z() > bMax.Z() || bMin.Z() > aMax.Z() {
			return false
		}

		return true
	}

	checkSphereBox := func(posA, posB mgl32.Vec3, colA SphereCollider, colB BoxCollider) bool {
		bMin := posB.Sub(colB.Origin)
		bMax := bMin.Add(colB.Size)

		nearest := mgl32.Vec3{
			m32.Max(bMin.X(), m32.Min(posA.X(), bMax.X())),
			m32.Max(bMin.Y(), m32.Min(posA.Y(), bMax.Y())),
			m32.Max(bMin.Z(), m32.Min(posA.Z(), bMax.Z())),
		}

		dist := DistanceSquared(posA, nearest)
		return dist < (colA.Radius * colA.Radius)
	}

	switch col := rb.Collider.(type) {
	case SphereCollider:
		switch otherCol := other.Collider.(type) {
		case SphereCollider:
			if checkSphereSphere(pos, otherPos, col, otherCol) {
				rb.Collide(other)
			}
		case BoxCollider:
			if checkSphereBox(pos, otherPos, col, otherCol) {
				rb.Collide(other)
			}
		}
	case BoxCollider:
		switch otherCol := other.Collider.(type) {
		case SphereCollider:
			if checkSphereBox(otherPos, pos, otherCol, col) {
				rb.Collide(other)
			}
		case BoxCollider:
			if checkBoxBox(pos, otherPos, col, otherCol) {
				rb.Collide(other)
			}
		}
	}
}

// Collide resolves a collision between two rigid bodies
func (rb *RigidBody) Collide(other *RigidBody) {
	diff := rb.Parent.Transform().Position.Add(other.Parent.Transform().Position.Mul(-1.0))

	x := diff.Normalize()
	v1 := rb.Velocity
	x1 := x.Dot(v1)
	v1x := x.Mul(x1)
	v1y := v1.Add(v1x.Mul(-1.0))
	m1 := rb.Mass

	x = x.Mul(-1.0)
	v2 := other.Velocity
	x2 := x.Dot(v2)
	v2x := x.Mul(x2)
	v2y := v2.Add(v2x.Mul(-1.0))
	m2 := other.Mass

	if m1 == m32.MaxFloat32 {
		other.Velocity = other.Velocity.Mul(-1.0 * other.Elasticity)
	} else if m2 == m32.MaxFloat32 {
		rb.Velocity = rb.Velocity.Mul(-1.0 * rb.Elasticity)
	} else {
		rb.Velocity = v1x.Mul((m1 - m2) / (m1 + m2)).Add(v2x.Mul((2 * m2) / (m1 + m2)).Add(v1y)).Mul(rb.Elasticity)
		other.Velocity = v1x.Mul((2 * m1) / (m1 + m2)).Add(v2x.Mul((m2 - m1) / (m1 + m2)).Add(v2y)).Mul(other.Elasticity)
	}
}
