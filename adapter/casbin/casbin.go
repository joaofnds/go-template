package casbin

const (
	RequestSubjectIndex = iota
	RequestDomainIndex
	RequestObjectIndex
	RequestActionIndex

	PolicySubjectIndex = iota
	PolicyDomainIndex
	PolicyObjectIndex
	PolicyActionIndex
	PolicyEffectIndex
)
