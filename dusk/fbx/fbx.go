package fbx

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"path/filepath"
	"strings"

	"github.com/WhoBrokeTheBuild/GoDusk/dusk"
	"github.com/go-gl/mathgl/mgl32"
)

func init() {
	dusk.RegisterMeshFormat("fbx", []string{".fbx"}, Load)
}

const (
	magic = "Kaydara FBX Binary  \x00"

	typeBool   uint8 = 'C'
	typeShort        = 'Y'
	typeInt          = 'I'
	typeLong         = 'L'
	typeFloat        = 'F'
	typeDouble       = 'D'

	typeBoolArray   = 'b'
	typeIntArray    = 'i'
	typeLongArray   = 'l'
	typeFloatArray  = 'f'
	typeDoubleArray = 'd'

	typeString = 'S'
	typeRaw    = 'R'
)

type header struct {
	Magic   [21]byte
	_       [2]byte
	Version uint32
}

type node struct {
	Name  string
	Props []*prop
	Nodes []*node
}

type prop struct {
	Type  uint8
	Value interface{}
}

type reader struct {
	Buffer    *bytes.Buffer
	Header    *header
	EndOffset int
}

func readArray(r *reader, size int) (data []byte, err error) {
	var len uint32
	err = binary.Read(r.Buffer, binary.LittleEndian, &len)
	if err != nil {
		panic(err)
		return
	}

	var encoding uint32
	err = binary.Read(r.Buffer, binary.LittleEndian, &encoding)
	if err != nil {
		panic(err)
		return
	}

	var compLen uint32
	err = binary.Read(r.Buffer, binary.LittleEndian, &compLen)
	if err != nil {
		panic(err)
		return
	}

	if encoding == 1 {
		tmp := r.Buffer.Next(int(compLen))

		var z io.ReadCloser
		z, err = zlib.NewReader(bytes.NewReader(tmp))
		if err != nil {
			return
		}
		defer z.Close()

		data, err = ioutil.ReadAll(z)
		if err != nil {
			return
		}
	} else {
		data = r.Buffer.Next(int(len) * size)
	}

	return
}

func (p *prop) read(r *reader) (err error) {
	err = binary.Read(r.Buffer, binary.LittleEndian, &p.Type)
	if err != nil {
		return
	}

	if p.Type == 0 {
		return
	}

	switch p.Type {
	case typeBool:
		var tmp uint8
		err = binary.Read(r.Buffer, binary.LittleEndian, &tmp)
		if err != nil {
			return
		}
		p.Value = tmp
		return

	case typeShort:
		var tmp uint16
		err = binary.Read(r.Buffer, binary.LittleEndian, &tmp)
		if err != nil {
			return
		}
		p.Value = tmp
		return

	case typeInt:
		var tmp int32
		err = binary.Read(r.Buffer, binary.LittleEndian, &tmp)
		if err != nil {
			return
		}
		p.Value = tmp
		return

	case typeLong:
		var tmp int64
		err = binary.Read(r.Buffer, binary.LittleEndian, &tmp)
		if err != nil {
			return
		}
		p.Value = tmp
		return

	case typeFloat:
		var tmp float32
		err = binary.Read(r.Buffer, binary.LittleEndian, &tmp)
		if err != nil {
			return
		}
		p.Value = tmp
		return

	case typeDouble:
		var tmp float64
		err = binary.Read(r.Buffer, binary.LittleEndian, &tmp)
		if err != nil {
			return
		}
		p.Value = tmp
		return

	case typeString:
		fallthrough

	case typeRaw:
		var len uint32
		err = binary.Read(r.Buffer, binary.LittleEndian, &len)
		if err != nil {
			return
		}

		tmp := make([]byte, int(len))
		err = binary.Read(r.Buffer, binary.LittleEndian, tmp)
		if err != nil {
			return
		}

		if p.Type == typeString {
			p.Value = string(tmp)
		} else {
			p.Value = tmp
		}

		return
	}

	var size int

	switch p.Type {
	case typeBoolArray:
		size = 1
	case typeIntArray:
		size = 4
	case typeLongArray:
		size = 8
	case typeFloatArray:
		size = 4
	case typeDoubleArray:
		size = 8
	default:
		return
	}

	var data []byte
	data, err = readArray(r, size)
	if err != nil {
		return
	}

	sr := bytes.NewReader(data)

	switch p.Type {
	case typeBoolArray:
		tmp := make([]uint8, len(data)/size)
		for i := 0; i < len(tmp); i++ {
			binary.Read(sr, binary.LittleEndian, &tmp[i])
		}
		p.Value = tmp

	case typeIntArray:
		tmp := make([]int32, len(data)/size)
		for i := 0; i < len(tmp); i++ {
			binary.Read(sr, binary.LittleEndian, &tmp[i])
		}
		p.Value = tmp

	case typeLongArray:
		tmp := make([]int64, len(data)/size)
		for i := 0; i < len(tmp); i++ {
			binary.Read(sr, binary.LittleEndian, &tmp[i])
		}
		p.Value = tmp

	case typeFloatArray:
		tmp := make([]float32, len(data)/size)
		for i := 0; i < len(tmp); i++ {
			binary.Read(sr, binary.LittleEndian, &tmp[i])
		}
		p.Value = tmp

	case typeDoubleArray:
		tmp := make([]float64, len(data)/size)
		for i := 0; i < len(tmp); i++ {
			binary.Read(sr, binary.LittleEndian, &tmp[i])
		}
		p.Value = tmp
	}

	return
}

type nodeHeader32 struct {
	EndOffset uint32
	NumProps  uint32
	_         uint32
	NameLen   uint8
}

func (nh *nodeHeader32) GetEndOffset() int {
	return int(nh.EndOffset)
}

func (nh *nodeHeader32) GetNumProps() int {
	return int(nh.NumProps)
}

func (nh *nodeHeader32) GetNameLen() int {
	return int(nh.NameLen)
}

type nodeHeader64 struct {
	EndOffset uint64
	NumProps  uint64
	_         uint64
	NameLen   uint8
}

func (nh *nodeHeader64) GetEndOffset() int {
	return int(nh.EndOffset)
}

func (nh *nodeHeader64) GetNumProps() int {
	return int(nh.NumProps)
}

func (nh *nodeHeader64) GetNameLen() int {
	return int(nh.NameLen)
}

type nodeHeader interface {
	GetEndOffset() int
	GetNumProps() int
	GetNameLen() int
}

func (n *node) read(r *reader) (err error) {
	var h nodeHeader
	if r.Header.Version > 7500 {
		h = &nodeHeader64{}
	} else {
		h = &nodeHeader32{}
	}
	err = binary.Read(r.Buffer, binary.LittleEndian, h)
	if err != nil {
		return
	}

	if h.GetEndOffset() == 0 {
		return
	}

	n.Name = string(r.Buffer.Next(int(h.GetNameLen())))

	n.Props = []*prop{}
	n.Nodes = []*node{}

	for i := 0; i < int(h.GetNumProps()); i++ {
		np := &prop{}
		err = np.read(r)
		if err != nil {
			return
		}
		n.Props = append(n.Props, np)
	}

	for {
		off := r.EndOffset - r.Buffer.Len()
		if off == h.GetEndOffset() {
			break
		}
		nn := &node{}
		err = nn.read(r)
		if err != nil {
			return
		}
		n.Nodes = append(n.Nodes, nn)
	}

	return
}

func (n *node) findAll(name string) []*node {
	if n == nil {
		return nil
	}
	nodes := []*node{}
	for _, c := range n.Nodes {
		if c.Name == name {
			nodes = append(nodes, c)
		}
	}
	return nodes
}

func (n *node) findFirst(name string) *node {
	if n == nil {
		return nil
	}
	for _, c := range n.Nodes {
		if c.Name == name {
			return c
		}
	}
	return nil
}

func (n *node) findByID(id int64) *node {
	if n == nil {
		return nil
	}
	for _, c := range n.Nodes {
		if len(c.Props) > 0 {
			switch cID := c.Props[0].Value.(type) {
			case int64:
				if id == cID {
					return c
				}
			}
		}
	}
	return nil
}

type conn struct {
	Type string
	A    *node
	B    *node
	Bind string
}

func newConn(root, n *node) *conn {
	c := &conn{}

	if len(n.Props) < 3 {
		return nil
	}

	c.Type = n.Props[0].Value.(string)
	a := n.Props[1].Value.(int64)
	b := n.Props[2].Value.(int64)

	if c.Type == "OP" {
		if len(n.Props) >= 4 {
			c.Bind = n.Props[3].Value.(string)
		}
	}

	objNode := root.findFirst("Objects")
	if objNode == nil {
		return nil
	}

	c.A = objNode.findByID(a)
	c.B = objNode.findByID(b)

	return c
}

type propMap map[string]interface{}

func newPropMap(n *node) propMap {
	pm := propMap{}

	pNodes := n.findAll("P")
	for _, pNode := range pNodes {
		if len(pNode.Props) < 4 {
			dusk.Warnf("Not enough props in 'P' node")
			continue
		}
		pName := pNode.Props[0].Value.(string)
		pType := pNode.Props[1].Value.(string)

		switch pType {
		case "ColorRGB":
			fallthrough
		case "Vector3D":
			fallthrough
		case "Lcl Rotation":
			fallthrough
		case "Lcl Scaling":
			if len(pNode.Props) < 7 {
				dusk.Warnf("Not enough props in 'P' node")
				continue
			}
			pm[pName] = mgl32.Vec3{
				float32(pNode.Props[4].Value.(float64)),
				float32(pNode.Props[5].Value.(float64)),
				float32(pNode.Props[6].Value.(float64)),
			}
		}
	}
	return pm
}

// Load parses and returns the data from the file
func Load(filename string) (data []*dusk.MeshData, err error) {
	filename = filepath.Clean(filename)
	dir := filepath.Dir(filename)

	file, err := dusk.Load(filename)
	if err != nil {
		return
	}

	b := bytes.NewBuffer(file)

	r := &reader{
		Buffer:    b,
		Header:    &header{},
		EndOffset: b.Len(),
	}

	err = binary.Read(r.Buffer, binary.LittleEndian, r.Header)
	if err != nil {
		return
	}

	if string(r.Header.Magic[:]) != magic {
		return nil, fmt.Errorf("Invalid file format")
	}

	dusk.Verbosef("Binary FBX Version %d.%d", r.Header.Version/1000, (r.Header.Version%1000)/100)

	root := &node{}
	for r.Buffer.Len() > 0 {
		n := &node{}
		err = n.read(r)
		if err == io.EOF {
			err = nil
			root.Nodes = append(root.Nodes, n)
			break
		}
		if err != nil {
			return
		}
		root.Nodes = append(root.Nodes, n)
	}

	var ambient mgl32.Vec4
	var diffuse mgl32.Vec4
	var specular mgl32.Vec4

	defNode := root.findFirst("Definitions")
	if defNode != nil {
		objTypeNodes := defNode.findAll("ObjectType")
		for _, otNode := range objTypeNodes {
			if len(otNode.Props) == 0 {
				continue
			}
			t := otNode.Props[0].Value.(string)
			if t == "Geometry" || t == "Model" || t == "Material" {
				propTempNode := otNode.findFirst("PropertyTemplate")
				if propTempNode != nil {
					p70Node := propTempNode.findFirst("Properties70")
					if p70Node != nil {
						pMap := newPropMap(p70Node)

						if value, found := pMap["Color"].(mgl32.Vec3); found {
							diffuse = mgl32.Vec4{value[0], value[1], value[2], 1.0}
						}

						if value, found := pMap["AmbientColor"].(mgl32.Vec3); found {
							diffuse = mgl32.Vec4{value[0], value[1], value[2], 1.0}
						}
						if value, found := pMap["DiffuseColor"].(mgl32.Vec3); found {
							diffuse = mgl32.Vec4{value[0], value[1], value[2], 1.0}
						}
						if value, found := pMap["SpecularColor"].(mgl32.Vec3); found {
							specular = mgl32.Vec4{value[0], value[1], value[2], 1.0}
						}
					}
				}
			}
		}
	}

	conns := []*conn{}
	connNode := root.findFirst("Connections")
	if connNode != nil {
		cNodes := connNode.findAll("C")
		for _, cNode := range cNodes {
			c := newConn(root, cNode)
			if c != nil {
				conns = append(conns, c)
			}
		}
	}

	objNode := root.findFirst("Objects")
	if objNode == nil {
		return nil, fmt.Errorf("FBX has no 'Objects' node")
	}

	data = []*dusk.MeshData{}
	modelNodes := objNode.findAll("Model")
	for _, modelNode := range modelNodes {
		var geomNode *node
		var matNode *node

		rotation := mgl32.Vec3{0, 0, 0}
		scale := mgl32.Vec3{1, 1, 1}

		for _, c := range conns {
			if c.A == nil {
				continue
			}

			if c.B == modelNode {
				switch c.A.Name {
				case "Geometry":
					geomNode = c.A
				case "Material":
					matNode = c.A
				}
			}
		}

		if geomNode == nil {
			continue
		}

		var pMap propMap
		p70Node := modelNode.findFirst("Properties70")
		if p70Node != nil {
			pMap = newPropMap(p70Node)
		}

		if value, found := pMap["Lcl Rotation"].(mgl32.Vec3); found {
			rotation = value
		}

		if value, found := pMap["Lcl Scaling"].(mgl32.Vec3); found {
			scale = value.Mul(1.0 / 100.0)
		}

		var ambientTexNode *node
		var diffuseTexNode *node
		var specularTexNode *node
		var bumpTexNode *node

		for _, c := range conns {
			if c.A == nil {
				continue
			}

			if c.B == matNode {
				if c.A.Name == "Texture" {
					switch c.Bind {
					case "AmbientColor":
						ambientTexNode = c.A
					case "DiffuseColor":
						diffuseTexNode = c.A
					case "Specular":
						specularTexNode = c.A
					case "Bump":
						bumpTexNode = c.A
					}
				}
			}
		}

		matData := &dusk.MaterialData{
			Ambient:  ambient,
			Diffuse:  diffuse,
			Specular: specular,
		}

		for _, c := range conns {
			if c.A == nil {
				continue
			}

			if c.A.Name == "Video" {
				fileNode := c.A.findFirst("Filename")
				relFileNode := c.A.findFirst("RelativeFilename")
				if fileNode == nil && relFileNode == nil {
					continue
				}

				f := ""
				if relFileNode != nil && len(relFileNode.Props) > 0 {
					f = strings.Replace(relFileNode.Props[0].Value.(string), "\\", "/", -1)
				}
				// Empty, or Absolute Path (Linux/Windows)
				if f == "" || f[0] == '/' || f[1:3] == ":/" {
					if fileNode != nil && len(fileNode.Props) > 0 {
						f = strings.Replace(fileNode.Props[0].Value.(string), "\\", "/", -1)
					}
				}

				if f == "" {
					continue
				}

				f = filepath.Join(dir, filepath.Clean(f))

				if c.B == ambientTexNode {
					matData.AmbientMap = f
				} else if c.B == diffuseTexNode {
					matData.DiffuseMap = f
				} else if c.B == specularTexNode {
					matData.SpecularMap = f
				} else if c.B == bumpTexNode {
					matData.BumpMap = f
				}
			}
		}

		vertInds := []int32{}
		verts := []float64{}

		normInds := []int32{}
		norms := []float64{}
		normMapping := ""
		normReference := ""

		txcdInds := []int32{}
		txcds := []float64{}
		txcdMapping := ""
		txcdReference := ""

		vertNode := geomNode.findFirst("Vertices")
		if vertNode == nil || len(vertNode.Props) == 0 {
			dusk.Warnf("No 'Vertices' in 'Geometry' node")
			continue
		}

		verts = vertNode.Props[0].Value.([]float64)
		if len(verts) == 0 {
			dusk.Warnf("No vertices")
		}

		indexNode := geomNode.findFirst("PolygonVertexIndex")
		if indexNode == nil || len(indexNode.Props) == 0 {
			dusk.Warnf("No 'PolygonVertexIndex' in 'Geometry' node")
			continue
		}

		vertInds = indexNode.Props[0].Value.([]int32)
		if len(verts) == 0 {
			dusk.Warnf("No vertInds")
		}

		normLayerNode := geomNode.findFirst("LayerElementNormal")
		if normLayerNode != nil {
			normMapNode := normLayerNode.findFirst("MappingInformationType")
			if normMapNode == nil || len(normMapNode.Props) == 0 {
				dusk.Warnf("No 'MappingInformationType' in 'LayerElementNormal' node")
			} else {
				normMapping = normMapNode.Props[0].Value.(string)
			}

			normRefNode := normLayerNode.findFirst("ReferenceInformationType")
			if normRefNode == nil || len(normRefNode.Props) == 0 {
				dusk.Warnf("No 'ReferenceInformationType' in 'LayerElementNormal' node")
			} else {
				normReference = normRefNode.Props[0].Value.(string)
			}

			if normReference == "IndexToDirect" {
				normIndNode := normLayerNode.findFirst("NormalsIndex")
				if normIndNode == nil || len(normIndNode.Props) == 0 {
					dusk.Warnf("No 'NormalsIndex' in 'LayerElementNormal' node")
				} else {
					normInds = normIndNode.Props[0].Value.([]int32)
				}
			}

			normNode := normLayerNode.findFirst("Normals")
			if normNode == nil || len(normNode.Props) == 0 {
				dusk.Warnf("No 'Normals' in 'LayerElementNormal' node")
			} else {
				norms = normNode.Props[0].Value.([]float64)
			}
		}

		txcdLayerNode := geomNode.findFirst("LayerElementUV")
		if txcdLayerNode != nil {
			txcdMapNode := txcdLayerNode.findFirst("MappingInformationType")
			if txcdMapNode == nil || len(txcdMapNode.Props) == 0 {
				dusk.Warnf("No 'MappingInformationType' in 'LayerElementUV' node")
			} else {
				txcdMapping = txcdMapNode.Props[0].Value.(string)
			}

			txcdRefNode := txcdLayerNode.findFirst("ReferenceInformationType")
			if txcdRefNode == nil || len(txcdRefNode.Props) == 0 {
				dusk.Warnf("No 'ReferenceInformationType' in 'LayerElementUV' node")
			} else {
				txcdReference = txcdRefNode.Props[0].Value.(string)
			}

			if txcdReference == "IndexToDirect" {
				txcdIndNode := txcdLayerNode.findFirst("UVIndex")
				if txcdIndNode == nil || len(txcdIndNode.Props) == 0 {
					dusk.Warnf("No 'UVIndex' in 'LayerElementUV' node")
				} else {
					txcdInds = txcdIndNode.Props[0].Value.([]int32)
				}
			}

			txcdNode := txcdLayerNode.findFirst("UV")
			if txcdNode == nil || len(txcdNode.Props) == 0 {
				dusk.Warnf("No 'UV' in 'LayerElementUV' node")
			} else {
				txcds = txcdNode.Props[0].Value.([]float64)
			}
		}

		mat, err := dusk.NewMaterialFromData(matData)
		if err != nil {
			return nil, err
		}

		_ = normInds
		_ = txcdInds

		d := &dusk.MeshData{
			Vertices:  []mgl32.Vec3{},
			Normals:   []mgl32.Vec3{},
			TexCoords: []mgl32.Vec2{},
			Material:  mat,
		}
		inds := []int{}
		for i := 0; i < len(vertInds); i += len(inds) {
			inds = []int{}
			for j := i; j < len(vertInds); j++ {
				ind := vertInds[j]
				inds = append(inds, int(ind))
				if ind < 0 {
					break
				}
			}
			inds[len(inds)-1] ^= -1

			order := []int{0, 1, 2}
			if len(inds) == 4 {
				order = []int{0, 1, 2, 2, 3, 0}
			}

			for _, o := range order {
				ind := inds[o] * 3
				v := mgl32.Vec3{
					float32(verts[ind+0]),
					float32(verts[ind+1]),
					float32(verts[ind+2]),
				}
				v = mgl32.TransformCoordinate(v, mgl32.HomogRotate3DX(rotation[0]*(-180.0/math.Pi)))
				v = mgl32.TransformCoordinate(v, mgl32.HomogRotate3DY(rotation[1]*(-180.0/math.Pi)))
				v = mgl32.TransformCoordinate(v, mgl32.HomogRotate3DZ(rotation[2]*(-180.0/math.Pi)))
				v = mgl32.TransformCoordinate(v, mgl32.Scale3D(scale[0], scale[1], scale[2]))
				d.Vertices = append(d.Vertices, v)
			}

			if len(norms) > 0 {
				if normMapping == "ByPolygonVertex" {
					if normReference == "Direct" {
						for _, o := range order {
							ind := (i + o) * 3
							n := mgl32.Vec3{
								float32(norms[ind+0]),
								float32(norms[ind+1]),
								float32(norms[ind+2]),
							}
							n = mgl32.TransformNormal(n, mgl32.HomogRotate3DX(rotation[0]*(-180.0/math.Pi)))
							n = mgl32.TransformNormal(n, mgl32.HomogRotate3DY(rotation[1]*(-180.0/math.Pi)))
							n = mgl32.TransformNormal(n, mgl32.HomogRotate3DZ(rotation[2]*(-180.0/math.Pi)))
							d.Normals = append(d.Normals, n)
						}
					} else if normReference == "IndexToDirect" {
						for _, o := range order {
							ind := inds[i+o] * 3
							n := mgl32.Vec3{
								float32(norms[ind+0]),
								float32(norms[ind+1]),
								float32(norms[ind+2]),
							}
							n = mgl32.TransformNormal(n, mgl32.HomogRotate3DX(rotation[0]*(-180.0/math.Pi)))
							n = mgl32.TransformNormal(n, mgl32.HomogRotate3DY(rotation[1]*(-180.0/math.Pi)))
							n = mgl32.TransformNormal(n, mgl32.HomogRotate3DZ(rotation[2]*(-180.0/math.Pi)))
							d.Normals = append(d.Normals, n)
						}
					}
				} else if normMapping == "ByVertex" {
					if normReference == "Direct" {
						// TODO
					} else if normReference == "IndexToDirect" {
						// TODO
					}
				}
			}

			if len(txcds) > 0 {
				if txcdMapping == "ByPolygonVertex" {
					if txcdReference == "Direct" {
						for _, o := range order {
							ind := (i + o) * 2
							d.TexCoords = append(d.TexCoords, mgl32.Vec2{
								float32(txcds[ind+0]),
								float32(txcds[ind+1]),
							})
						}
					} else if txcdReference == "IndexToDirect" {
						for _, o := range order {
							ind := txcdInds[i+o] * 2
							d.TexCoords = append(d.TexCoords, mgl32.Vec2{
								float32(txcds[ind+0]),
								float32(txcds[ind+1]),
							})
						}
					}
				} else if txcdMapping == "ByVertex" {
					if txcdReference == "Direct" {
						// TODO
					} else if txcdReference == "IndexToDirect" {
						// TODO
					}
				}
			}
		}

		data = append(data, d)
	}

	return
}
