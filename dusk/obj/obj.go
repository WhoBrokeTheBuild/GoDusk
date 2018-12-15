package obj

import (
	"bytes"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/go-gl/mathgl/mgl32"

	"github.com/WhoBrokeTheBuild/GoDusk/dusk"
)

func init() {
	dusk.RegisterMeshFormat("obj", []string{".obj"}, Load)
}

// Load parses and returns the data from the file
func Load(filename string) ([]*dusk.MeshData, error) {
	filename = filepath.Clean(filename)
	dir := filepath.Dir(filename)

	data := []*dusk.MeshData{}
	materials := map[string]*dusk.MaterialData{}

	file, err := dusk.Load(filename)
	if err != nil {
		return nil, err
	}

	if file[len(file)-1] != '\n' {
		file = append(file, '\n')
	}

	buf := bytes.NewBuffer(file)

	var o *dusk.MeshData
	var x, y, z, u, v float32
	var f [4][3]int

	var hasNorm bool
	var hasTxcd bool

	verts := []mgl32.Vec3{}
	norms := []mgl32.Vec3{}
	txcds := []mgl32.Vec2{}

	for {
		bytes, err := buf.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return nil, err
			}
		}

		line := string(bytes)

		if len(line) < 2 || line[0] == '#' || line[0] == '\n' {
			continue
		}

		if line[0] == 'v' {
			if line[1] == 'n' {
				fmt.Sscanf(line[3:], "%f %f %f", &x, &y, &z)
				norms = append(norms, mgl32.Vec3{x, y, z})
			} else if line[1] == 't' {
				fmt.Sscanf(line[3:], "%f %f", &u, &v)
				txcds = append(txcds, mgl32.Vec2{u, v})
			} else {
				fmt.Sscanf(line[2:], "%f %f %f", &x, &y, &z)
				verts = append(verts, mgl32.Vec3{x, y, z})
			}
		} else if line[0] == 'f' {
			if o == nil {
				o = &dusk.MeshData{}
				data = append(data, o)
			}

			polySize := 3
			hasNorm = false
			hasTxcd = false
			if strings.Contains(line, "//") {
				hasNorm = true
				polySize, _ = fmt.Sscanf(line[2:],
					"%d//%d %d//%d %d//%d %d//%d",
					&f[0][0], &f[0][2],
					&f[1][0], &f[1][2],
					&f[2][0], &f[2][2],
					&f[3][0], &f[3][2])
				polySize /= 2
			} else {
				sc := strings.Count(line, "/")
				if sc == 3 || sc == 4 {
					hasTxcd = true
					polySize, _ = fmt.Sscanf(line[2:],
						"%d/%d %d/%d %d/%d %d/%d",
						&f[0][0], &f[0][1],
						&f[1][0], &f[1][1],
						&f[2][0], &f[2][1],
						&f[2][0], &f[2][1])
					polySize /= 2
				} else {
					hasNorm = true
					hasTxcd = true
					polySize, _ = fmt.Sscanf(line[2:],
						"%d/%d/%d %d/%d/%d %d/%d/%d %d/%d/%d",
						&f[0][0], &f[0][1], &f[0][2],
						&f[1][0], &f[1][1], &f[1][2],
						&f[2][0], &f[2][1], &f[2][2],
						&f[3][0], &f[3][1], &f[3][2])
					polySize /= 3
				}
			}
			// TODO: Handle `f %d %d %d`

			inds := []int{0, 1, 2}
			if polySize == 4 {
				inds = []int{0, 1, 2, 2, 3, 0}
			}
			for _, i := range inds {
				if f[i][0] < 0 {
					f[i][0] += len(verts) + 1
				}
				if f[i][1] < 0 {
					f[i][1] += len(txcds) + 1
				}
				if f[i][2] < 0 {
					f[i][2] += len(norms) + 1
				}

				o.Vertices = append(o.Vertices, verts[f[i][0]-1])

				if hasTxcd {
					o.TexCoords = append(o.TexCoords, txcds[f[i][1]-1])
				}

				if hasNorm {
					o.Normals = append(o.Normals, norms[f[i][2]-1])
				}
			}
		} else if line[0] == 'o' {
			name := strings.TrimSpace(line[2:])
			dusk.Verbosef("Processing Object [%v]", name)
			if o != nil && o.Name == "" {
				o.Name = name
			} else {
				o = &dusk.MeshData{
					Name:      name,
					Vertices:  []mgl32.Vec3{},
					Normals:   []mgl32.Vec3{},
					TexCoords: []mgl32.Vec2{},
				}
				data = append(data, o)
			}
		} else if strings.HasPrefix(line, "mtllib") {
			tmp, err := readMaterial(filepath.Join(dir, strings.TrimSpace(line[7:])))
			if err != nil {
				return nil, err
			}
			for k, v := range tmp {
				materials[k] = v
			}
		} else if strings.HasPrefix(line, "usemtl") {
			if o == nil {
				o = &dusk.MeshData{}
				data = append(data, o)
			}

			name := strings.TrimSpace(line[7:])
			if m, ok := materials[name]; ok {
				o.Material, err = dusk.NewMaterialFromData(m)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	return data, nil
}

func readMaterial(filename string) (map[string]*dusk.MaterialData, error) {
	filename = filepath.Clean(filename)
	dir := filepath.Dir(filename)

	materials := map[string]*dusk.MaterialData{}

	file, err := dusk.Load(filename)
	if err != nil {
		return nil, err
	}

	if file[len(file)-1] != '\n' {
		file = append(file, '\n')
	}

	buf := bytes.NewBuffer(file)

	var m *dusk.MaterialData

	for {
		bytes, err := buf.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return nil, err
			}
		}

		line := string(bytes)

		if len(line) < 2 || line[0] == '#' || line[0] == '\n' {
			continue
		}

		if strings.HasPrefix(line, "newmtl") { // newmtl
			name := strings.TrimSpace(line[7:])
			m = &dusk.MaterialData{
				Ambient:  mgl32.Vec4{0, 0, 0, 1},
				Diffuse:  mgl32.Vec4{0, 0, 0, 1},
				Specular: mgl32.Vec4{0, 0, 0, 1},
			}
			materials[name] = m
		} else if line[0] == 'K' {
			if line[1] == 'a' {
				if m != nil {
					// Ka
					fmt.Sscanf(line[3:], "%f %f %f", &m.Ambient[0], &m.Ambient[1], &m.Ambient[2])
				}
			} else if line[1] == 'd' {
				if m != nil {
					// Kd
					fmt.Sscanf(line[3:], "%f %f %f", &m.Diffuse[0], &m.Diffuse[1], &m.Diffuse[2])
				}
			} else if line[1] == 's' {
				if m != nil {
					// Ks
					fmt.Sscanf(line[3:], "%f %f %f", &m.Specular[0], &m.Specular[1], &m.Specular[2])
				}
			}
		} else if line[0] == 'N' && line[1] == 's' {
			// Ns
			//fmt.Sscanf(line[3:], "%f", &m.Shininess)
		} else if strings.HasPrefix(line, "map_K") {
			if line[5] == 'a' {
				// map_Ka
				m.AmbientMap = filepath.Join(dir, strings.TrimSpace(line[7:]))
			} else if line[5] == 'd' {
				// map_Kd
				m.DiffuseMap = filepath.Join(dir, strings.TrimSpace(line[7:]))
			} else if line[5] == 's' {
				// map_Ks
				m.SpecularMap = filepath.Join(dir, strings.TrimSpace(line[7:]))
			}
		} else if strings.HasPrefix(line, "map_Ns") {
			// map_Ns
			//m.SpecularHighlightMap = filepath.Join(dir, strings.TrimSpace(line[7:]))
		} else if strings.HasPrefix(line, "bump") {
			// bump
			//m.BumpMap = filepath.Join(dir, strings.TrimSpace(line[5:]))
		} else if strings.HasPrefix(line, "map_bump") {
			// map_bump
			//m.BumpMap = filepath.Join(dir, strings.TrimSpace(line[7:]))
		} else if strings.HasPrefix(line, "disp") {
			// disp
			//m.DisplacementMap = filepath.Join(dir, strings.TrimSpace(line[5:]))
		} else if strings.HasPrefix(line, "refl") {
			// refl
			//m.ReflectionMap = filepath.Join(dir, strings.TrimSpace(line[5:]))
		} else if strings.HasPrefix(line, "map_d") {
			// map_d
			//m.AlphaMap = filepath.Join(dir, strings.TrimSpace(line[6:]))
		}
	}

	return materials, nil
}
