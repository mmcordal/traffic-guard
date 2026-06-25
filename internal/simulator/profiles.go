package simulator

type DomainProfile struct {
	Domain string
	Weight int
}

type IPProfile struct {
	IP      string
	Country string
	ASN     string
	Weight  int
}
