package pkg


var (
	namespaces = make(map[string][]string, 10)
)

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
		if end == 0 {
			return deployIndex[start:]
		}else{
			return deployIndex[start:end]
		}
	}else{
		return nil
	}
}