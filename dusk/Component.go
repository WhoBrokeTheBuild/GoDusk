package dusk

type IComponent interface {
	Init(IEntity)
	Delete()

	GetEntity() IEntity

	Update(*UpdateContext)
	Render(*RenderContext)
}

type Component struct {
	entity IEntity
}

func NewComponent(entity IEntity) *Component {
	c := &Component{}
	c.Init(entity)
	return c
}

func (c *Component) Init(entity IEntity) {
	c.entity = entity
}

func (c *Component) Delete() {
	c.entity = nil
}

func (c *Component) GetEntity() IEntity {
	return c.entity
}

func (c *Component) Update(ctx *UpdateContext) {}

func (c *Component) Render(ctx *RenderContext) {}
