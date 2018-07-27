package packfile

import (
	"bytes"
	"io"
	"math"

	"io/ioutil"

	. "gopkg.in/check.v1"
	"gopkg.in/src-d/go-billy.v4/osfs"
	fixtures "gopkg.in/src-d/go-git-fixtures.v3"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/format/idxfile"
	"gopkg.in/src-d/go-git.v4/storage/memory"
)

type PackfileSuite struct {
	fixtures.Suite
	p   *Packfile
	idx *idxfile.MemoryIndex
	f   *fixtures.Fixture
}

var _ = Suite(&PackfileSuite{})

func (s *PackfileSuite) TestGet(c *C) {
	for h := range expectedEntries {
		obj, err := s.p.Get(h)
		c.Assert(err, IsNil)
		c.Assert(obj, Not(IsNil))
		c.Assert(obj.Hash(), Equals, h)
	}

	_, err := s.p.Get(plumbing.ZeroHash)
	c.Assert(err, Equals, plumbing.ErrObjectNotFound)
}

func (s *PackfileSuite) TestGetByOffset(c *C) {
	for h, o := range expectedEntries {
		obj, err := s.p.GetByOffset(o)
		c.Assert(err, IsNil)
		c.Assert(obj, Not(IsNil))
		c.Assert(obj.Hash(), Equals, h)
	}

	_, err := s.p.GetByOffset(math.MaxInt64)
	c.Assert(err, Equals, io.EOF)
}

func (s *PackfileSuite) TestID(c *C) {
	id, err := s.p.ID()
	c.Assert(err, IsNil)
	c.Assert(id, Equals, s.f.PackfileHash)
}

func (s *PackfileSuite) TestGetAll(c *C) {
	iter, err := s.p.GetAll()
	c.Assert(err, IsNil)

	var objects int
	for {
		o, err := iter.Next()
		if err == io.EOF {
			break
		}
		c.Assert(err, IsNil)

		objects++
		_, ok := expectedEntries[o.Hash()]
		c.Assert(ok, Equals, true)
	}

	c.Assert(objects, Equals, len(expectedEntries))
}

var expectedEntries = map[plumbing.Hash]int64{
	plumbing.NewHash("1669dce138d9b841a518c64b10914d88f5e488ea"): 615,
	plumbing.NewHash("32858aad3c383ed1ff0a0f9bdf231d54a00c9e88"): 1524,
	plumbing.NewHash("35e85108805c84807bc66a02d91535e1e24b38b9"): 1063,
	plumbing.NewHash("49c6bb89b17060d7b4deacb7b338fcc6ea2352a9"): 78882,
	plumbing.NewHash("4d081c50e250fa32ea8b1313cf8bb7c2ad7627fd"): 84688,
	plumbing.NewHash("586af567d0bb5e771e49bdd9434f5e0fb76d25fa"): 84559,
	plumbing.NewHash("5a877e6a906a2743ad6e45d99c1793642aaf8eda"): 84479,
	plumbing.NewHash("6ecf0ef2c2dffb796033e5a02219af86ec6584e5"): 186,
	plumbing.NewHash("7e59600739c96546163833214c36459e324bad0a"): 84653,
	plumbing.NewHash("880cd14280f4b9b6ed3986d6671f907d7cc2a198"): 78050,
	plumbing.NewHash("8dcef98b1d52143e1e2dbc458ffe38f925786bf2"): 84741,
	plumbing.NewHash("918c48b83bd081e863dbe1b80f8998f058cd8294"): 286,
	plumbing.NewHash("9a48f23120e880dfbe41f7c9b7b708e9ee62a492"): 80998,
	plumbing.NewHash("9dea2395f5403188298c1dabe8bdafe562c491e3"): 84032,
	plumbing.NewHash("a39771a7651f97faf5c72e08224d857fc35133db"): 84430,
	plumbing.NewHash("a5b8b09e2f8fcb0bb99d3ccb0958157b40890d69"): 838,
	plumbing.NewHash("a8d315b2b1c615d43042c3a62402b8a54288cf5c"): 84375,
	plumbing.NewHash("aa9b383c260e1d05fbbf6b30a02914555e20c725"): 84760,
	plumbing.NewHash("af2d6a6954d532f8ffb47615169c8fdf9d383a1a"): 449,
	plumbing.NewHash("b029517f6300c2da0f4b651b8642506cd6aaf45d"): 1392,
	plumbing.NewHash("b8e471f58bcbca63b07bda20e428190409c2db47"): 1230,
	plumbing.NewHash("c192bd6a24ea1ab01d78686e417c8bdc7c3d197f"): 1713,
	plumbing.NewHash("c2d30fa8ef288618f65f6eed6e168e0d514886f4"): 84725,
	plumbing.NewHash("c8f1d8c61f9da76f4cb49fd86322b6e685dba956"): 80725,
	plumbing.NewHash("cf4aa3b38974fb7d81f367c0830f7d78d65ab86b"): 84608,
	plumbing.NewHash("d3ff53e0564a9f87d8e84b6e28e5060e517008aa"): 1685,
	plumbing.NewHash("d5c0f4ab811897cadf03aec358ae60d21f91c50d"): 2351,
	plumbing.NewHash("dbd3641b371024f44d0e469a9c8f5457b0660de1"): 84115,
	plumbing.NewHash("e8d3ffab552895c19b9fcf7aa264d277cde33881"): 12,
	plumbing.NewHash("eba74343e2f15d62adedfd8c883ee0262b5c8021"): 84708,
	plumbing.NewHash("fb72698cab7617ac416264415f13224dfd7a165e"): 84671,
}

func (s *PackfileSuite) TestContent(c *C) {
	storer := memory.NewObjectStorage()
	decoder, err := NewDecoder(NewScanner(s.f.Packfile()), storer)
	c.Assert(err, IsNil)

	_, err = decoder.Decode()
	c.Assert(err, IsNil)

	iter, err := s.p.GetAll()
	c.Assert(err, IsNil)

	for {
		o, err := iter.Next()
		if err == io.EOF {
			break
		}
		c.Assert(err, IsNil)

		o2, err := storer.EncodedObject(plumbing.AnyObject, o.Hash())
		c.Assert(err, IsNil)

		c.Assert(o.Type(), Equals, o2.Type())
		c.Assert(o.Size(), Equals, o2.Size())

		r, err := o.Reader()
		c.Assert(err, IsNil)

		c1, err := ioutil.ReadAll(r)
		c.Assert(err, IsNil)
		c.Assert(r.Close(), IsNil)

		r, err = o2.Reader()
		c.Assert(err, IsNil)

		c2, err := ioutil.ReadAll(r)
		c.Assert(err, IsNil)
		c.Assert(r.Close(), IsNil)

		c.Assert(bytes.Compare(c1, c2), Equals, 0)
	}
}

func (s *PackfileSuite) SetUpTest(c *C) {
	s.f = fixtures.Basic().One()

	f, err := osfs.New("/").Open(s.f.Packfile().Name())
	c.Assert(err, IsNil)

	s.idx = idxfile.NewMemoryIndex()
	c.Assert(idxfile.NewDecoder(s.f.Idx()).Decode(s.idx), IsNil)

	s.p = NewPackfile(s.idx, f)
}

func (s *PackfileSuite) TearDownTest(c *C) {
	c.Assert(s.p.Close(), IsNil)
}