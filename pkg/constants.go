package pkg

const (
	kubectlPrefix = "kubectl.kubernetes.io/"

	// LastAppliedConfigAnnotation is the annotation used to store the previous
	// configuration of a resource for use in a three way diff by UpdateApplyAnnotation.
	LastAppliedConfigAnnotation = kubectlPrefix + "last-applied-configuration"
)
