package pkg

import "sync"



type DeployIndexStore interface {
	Add(key, namespace string)
	Exists(key, namespace string) bool
	Delete(index, namespace string) bool
	GetAll(namespace string) []string
	Len(namespace string) int
	GetLimit(start, end int, namespace string) []string
}

type deployIndexMap struct {
	lock  sync.RWMutex
	items map[string][]string
}


func NewDeployIndexStore() DeployIndexStore {
	return &deployIndexMap{
		items:    make(map[string][]string, 10),
	}
}


// create index for k8s Deploy
// Current use for paging
// param index eg : Deploy.Name
func (d *deployIndexMap) Add(index, namespace string)  {
	if d.Exists(index, namespace) == false{
		d.lock.Lock()
		defer d.lock.Unlock()
		d.items[namespace] = append(d.items[namespace], index)
	}
}


func (d *deployIndexMap) Exists(index, namespace string) bool {
	d.lock.Lock()
	defer d.lock.Unlock()

	var deployIndex []string
	if _, ok := d.items[namespace]; ok{
		deployIndex = d.items[namespace]
	}else{
		return false
	}
	for _,v := range deployIndex {
		if v == index{
			return true
		}
	}

	return false
}

// Del index
//will rearrange after del
func (d *deployIndexMap) Delete(index, namespace string) bool {
	d.lock.Lock()
	defer d.lock.Unlock()

	var deployIndex []string
	if _, ok := d.items[namespace]; ok{
		deployIndex = d.items[namespace]
	}else{
		return false
	}
	for k,v := range deployIndex {
		if v == index{
			deployIndex = append(deployIndex[:k], deployIndex[k+1:]...)
			d.items[namespace] = deployIndex
			return true
		}
	}

	return false
}

// return index length
func (d *deployIndexMap) Len(namespace string) int {
	d.lock.RLock()
	defer d.lock.RUnlock()

	var deployIndex []string
	if _, ok := d.items[namespace]; ok{
		deployIndex = d.items[namespace]
		return len(deployIndex)
	}else{
		return 0
	}
}


func (d *deployIndexMap) GetAll(namespace string) []string {
	d.lock.RLock()
	defer d.lock.RUnlock()

	var deployIndex []string
	if _, ok := d.items[namespace]; ok{
		deployIndex = d.items[namespace]
		return deployIndex
	}else{
		return nil
	}
}


func (d *deployIndexMap) GetLimit(start, end int, namespace string) []string {
	d.lock.RLock()
	defer d.lock.RUnlock()

	var deployIndex []string
	if _, ok := d.items[namespace]; ok{
		deployIndex = d.items[namespace]
		if start > len(deployIndex){
			if len(deployIndex) > 10 {
				return deployIndex[0:10]
			}else{
				return deployIndex[0:]
			}
		} else if end == 0 {
			return deployIndex[start:]
		}else{
			return deployIndex[start:end]
		}
	}else{
		return nil
	}
}