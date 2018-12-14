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

var ignore = []string{
	"FBXHeaderExtension",
	"FileId",
	"CreationTime",
	"Creator",
	"GlobalSettings",
	"Documents",
	"Takes",
}

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
		skip := false
		for _, name := range ignore {
			if name == n.Name {
				skip = true
			}
		}
		if !skip {
			n.Nodes = append(n.Nodes, nn)
		}
	}

	return
}

func (n *node) find(name string) []*node {
	nodes := []*node{}
	for _, c := range n.Nodes {
		if c.Name == name {
			nodes = append(nodes, c)
		}
	}
	return nodes
}

func (n *node) findFirst(name string) *node {
	for _, c := range n.Nodes {
		if c.Name == name {
			return c
		}
	}
	return nil
}

func (n *node) findByID(id int64) *node {
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

// Load parses and returns the data from the file
func Load(filename string) (data []*dusk.MeshData, err error) {
	filename = filepath.Clean(filename)
	dir := filepath.Dir(filename)

	//materials := map[string]*dusk.MaterialData{}

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
		skip := false
		for _, name := range ignore {
			if name == n.Name {
				skip = true
			}
		}
		if !skip {
			root.Nodes = append(root.Nodes, n)
		}
	}

	//j, _ := json.MarshalIndent(&root, "", "  ")
	//ioutil.WriteFile("test.json", j, os.ModePerm)

	connNode := root.findFirst("Connections")
	if connNode == nil {
		return nil, fmt.Errorf("FBX has no 'Connections' node")
	}

	cNodes := connNode.find("C")

	objNode := root.findFirst("Objects")
	if objNode == nil {
		return nil, fmt.Errorf("FBX has no 'Objects' node")
	}

	data = []*dusk.MeshData{}
	rotation := mgl32.Vec3{0, 0, 0}
	scale := mgl32.Vec3{1, 1, 1}

	modelNodes := objNode.find("Model")
	materialNodes := objNode.find("Material")
	videoNodes := objNode.find("Video")

	geomNodes := objNode.find("Geometry")
	for _, geomNode := range geomNodes {
		geomID := geomNode.Props[0].Value.(int64)

		conns := []int64{}

		for _, cNode := range cNodes {
			a := cNode.Props[1].Value.(int64)
			b := cNode.Props[2].Value.(int64)
			if geomID == a {
				conns = append(conns, b)
			}
		}

		modelID := int64(0)
		for _, modelNode := range modelNodes {
			tmp := modelNode.Props[0].Value.(int64)
			for _, id := range conns {
				if id == tmp {
					modelID = tmp
					propNode := modelNode.findFirst("Properties70")
					if propNode != nil {
						pNodes := propNode.find("P")
						for _, pNode := range pNodes {
							pName := pNode.Props[0].Value.(string)
							if pName == "Lcl Rotation" {
								rotation = mgl32.Vec3{
									float32(pNode.Props[4].Value.(float64)),
									float32(pNode.Props[5].Value.(float64)),
									float32(pNode.Props[6].Value.(float64)),
								}
							} else if pName == "Lcl Scaling" {
								scale = mgl32.Vec3{
									float32(pNode.Props[4].Value.(float64) / 100.0),
									float32(pNode.Props[5].Value.(float64) / 100.0),
									float32(pNode.Props[6].Value.(float64) / 100.0),
								}
							}
						}
					}
					break
				}
			}
		}

		conns = []int64{}
		for _, cNode := range cNodes {
			a := cNode.Props[1].Value.(int64)
			b := cNode.Props[2].Value.(int64)
			if modelID == b {
				conns = append(conns, a)
			}
		}

		matID := int64(0)
		for _, matNode := range materialNodes {
			tmp := matNode.Props[0].Value.(int64)
			for _, id := range conns {
				if id == tmp {
					matID = tmp
					break
				}
			}
		}

		diffTexID := int64(0)
		bumpTexID := int64(0)

		conns = []int64{}
		for _, cNode := range cNodes {
			a := cNode.Props[1].Value.(int64)
			b := cNode.Props[2].Value.(int64)
			if matID == b {
				if len(cNode.Props) >= 4 {
					bind := cNode.Props[3].Value.(string)
					switch bind {
					case "DiffuseColor":
						diffTexID = a
					case "Bump":
						bumpTexID = a
					}
				}
			}
		}

		diffVideoID := int64(0)
		bumpVideoID := int64(0)

		for _, cNode := range cNodes {
			a := cNode.Props[1].Value.(int64)
			b := cNode.Props[2].Value.(int64)
			if b == diffTexID {
				diffVideoID = a
			} else if b == bumpTexID {
				bumpVideoID = a
			}
		}

		diffFilename := ""
		bumpFilename := ""

		for _, vidNode := range videoNodes {
			tmp := vidNode.Props[0].Value.(int64)
			if tmp == diffVideoID {
				filenameNode := vidNode.findFirst("RelativeFilename")
				diffFilename = filenameNode.Props[0].Value.(string)
			} else if tmp == bumpVideoID {
				filenameNode := vidNode.findFirst("RelativeFilename")
				bumpFilename = filenameNode.Props[0].Value.(string)
			}
		}

		if diffFilename != "" {
			diffFilename = filepath.Join(dir, strings.Replace(diffFilename, "\\", "/", -1))
		}
		if bumpFilename != "" {
			bumpFilename = filepath.Join(dir, strings.Replace(bumpFilename, "\\", "/", -1))
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

		mat, err := dusk.NewMaterialFromData(&dusk.MaterialData{
			Diffuse:    mgl32.Vec4{0.8, 0.8, 0.8, 1.0},
			DiffuseMap: diffFilename,
			BumpMap:    bumpFilename,
		})
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
							ind := inds[o] * 3
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

		//dusk.Verbosef("%v", d.TexCoords)

		data = append(data, d)
	}

	return
}
