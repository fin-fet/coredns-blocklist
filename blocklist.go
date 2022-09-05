package blocklist

type Blocklist interface {
	Create()
	Add(domain string)
	Contains(domain string) bool
	Len() int
}

type BasicBlocklist struct {
	blockList map[string]struct{}
}

func NewBasicBlocklist() Blocklist {
	bbl := &BasicBlocklist{}
	bbl.Create()
	return bbl
}

func (bbl *BasicBlocklist) Create() {
	bbl.blockList = make(map[string]struct{})
}

func (bbl *BasicBlocklist) Add(name string) {
	bbl.blockList[name] = struct{}{}
}

func (bbl *BasicBlocklist) Len() int {
	return len(bbl.blockList)
}

func (bbl *BasicBlocklist) Contains(name string) bool {
	_, ok := bbl.blockList[name]
	return ok
}
