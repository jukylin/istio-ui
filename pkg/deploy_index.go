package pkg


var (
	namespaces = make(map[string][]string, 10)
)

// create index for k8s Deploy
// Current use for paging
// param index eg : Deploy.Name
func SetDeployIndex(index, namespace string) bool {

	namespaces[namespace] = append(namespaces[namespace], index)

	return true
}


func ExistsDeployIndex(index, namespace string) bool {
	var deployIndex []string
	if _, ok := namespaces[namespace]; ok{
		deployIndex = namespaces[namespace]
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
func DelDeployIndex(index, namespace string) bool {
	var deployIndex []string
	if _, ok := namespaces[namespace]; ok{
		deployIndex = namespaces[namespace]
	}else{
		return false
	}
	for k,v := range deployIndex {
		if v == index{
			deployIndex = append(deployIndex[:k], deployIndex[k+1:]...)
			namespaces[namespace] = deployIndex
			return true
		}
	}

	return false
}

// return index length
func DeployIndexLen(namespace string) int {
	var deployIndex []string
	if _, ok := namespaces[namespace]; ok{
		deployIndex = namespaces[namespace]
		return len(deployIndex)
	}else{
		return 0
	}
}


func GetAllDeployIndex(namespace string) []string {
	var deployIndex []string
	if _, ok := namespaces[namespace]; ok{
		deployIndex = namespaces[namespace]
		return deployIndex
	}else{
		return nil
	}
}


func GetDeployIndexLimit(start, end int, namespace string) []string {
	var deployIndex []string
	if _, ok := namespaces[namespace]; ok{
		deployIndex = namespaces[namespace]
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