/*
A direct translation of the XML described in RFC7489 Appendix C, representing an
aggregate DMARC report generated by a receiver.

I am including the XML specification for each type in order to make sure it's
all covered. I am, however, omitting the comments from the RFC. For more
information, see https://tools.ietf.org/html/rfc7489#appendix-C.
*/
package dmarc

import (
	"encoding/xml"
)

/*
Represents the range of time which the report covers. Specified in seconds since
epoch (UTC).

	<xs:complexType name="DateRangeType">
	  <xs:all>
	    <xs:element name="begin" type="xs:integer"/>
	    <xs:element name="end" type="xs:integer"/>
	  </xs:all>
	</xs:complexType>
*/
type DateRange struct {
	Begin int `xml:"begin"`
	End   int `xml:"end"`
}

/*
Metadata about the receiver generating the report, and the report itself.

	<xs:complexType name="ReportMetadataType">
	  <xs:sequence>
	    <xs:element name="org_name" type="xs:string"/>
	    <xs:element name="email" type="xs:string"/>
	    <xs:element name="extra_contact_info" type="xs:string"
			minOccurs="0"/>
	    <xs:element name="report_id" type="xs:string"/>
	    <xs:element name="date_range" type="DateRangeType"/>
	    <xs:element name="error" type="xs:string" minOccurs="0"
			maxOccurs="unbounded"/>
	  </xs:sequence>
	</xs:complexType>
*/
type Metadata struct {
	OrgName          string    `xml:"org_name"`
	Email            string    `xml:"email"`
	ExtraContactInfo string    `xml:"extra_contact_info,omitempty"`
	ReportId         string    `xml:"report_id"`
	DateRange        DateRange `xml:"date_range"`
	Error            []string  `xml:"error,omitempty"`
}

/*
Alignment mode. "r" for relaxed, or "s" for strict.

	<xs:simpleType name="AlignmentType">
	  <xs:restriction base="xs:string">
	    <xs:enumeration value="r"/>
	    <xs:enumeration value="s"/>
	  </xs:restriction>
	</xs:simpleType>
*/
type Alignment string

/*
The applied DMARC disposition mode.

	<xs:simpleType name="DispositionType">
	  <xs:restriction base="xs:string">
	    <xs:enumeration value="none"/>
	    <xs:enumeration value="quarantine"/>
	    <xs:enumeration value="reject"/>
	  </xs:restriction>
	</xs:simpleType>
*/
type Disposition string

/*
The DMARC policy applied to the messages.

	<xs:complexType name="PolicyPublishedType">
	  <xs:all>
	    <xs:element name="domain" type="xs:string"/>
	    <xs:element name="adkim" type="AlignmentType" minOccurs="0"/>
	    <xs:element name="aspf" type="AlignmentType" minOccurs="0"/>
	    <xs:element name="p" type="DispositionType"/>
	    <xs:element name="sp" type="DispositionType"/>
	    <xs:element name="pct" type="xs:integer"/>
	    <xs:element name="fo" type="xs:string"/>
	  </xs:all>
	</xs:complexType>
*/
type PolicyPublished struct {
	Domain string      `xml:"domain"`
	Adkim  Alignment   `xml:"adkim,omitempty"`
	Aspf   Alignment   `xml:"aspf,omitempty"`
	P      Disposition `xml:"p"`
	Sp     Disposition `xml:"sp"`
	Pct    int         `xml:"pct"`
	Fo     string      `xml:"fo,omitempty"`
}

/*
The result of the DMARC-aligned authentication. "pass" or "fail".

	<xs:simpleType name="DMARCResultType">
	  <xs:restriction base="xs:string">
	    <xs:enumeration value="pass"/>
	    <xs:enumeration value="fail"/>
	  </xs:restriction>
	</xs:simpleType>
*/
type DMARCResult string

/*
Reason for deviating from the Domain Owener's published policy. See RFC7489,
section 6.

	<xs:simpleType name="PolicyOverrideType">
	  <xs:restriction base="xs:string">
	    <xs:enumeration value="forwarded"/>
	    <xs:enumeration value="sampled_out"/>
	    <xs:enumeration value="trusted_forwarder"/>
	    <xs:enumeration value="mailing_list"/>
	    <xs:enumeration value="local_policy"/>
	    <xs:enumeration value="other"/>
	  </xs:restriction>
	</xs:simpleType>
*/
type PolicyOverrideType string

/*
Reason for deviating from the published policy, with optional comment.

	<xs:complexType name="PolicyOverrideReason">
	  <xs:all>
	    <xs:element name="type" type="PolicyOverrideType"/>
	    <xs:element name="comment" type="xs:string" minOccurs="0"/>
	  </xs:all>
	</xs:complexType>
*/
type PolicyOverrideReason struct {
	Type    PolicyOverrideType `xml:"type"`
	Comment string             `xml:"comment,omitempty"`
}

/*
The result after applying the DMARC policy.

	<xs:complexType name="PolicyEvaluatedType">
	  <xs:sequence>
	    <xs:element name="disposition" type="DispositionType"/>
	    <xs:element name="dkim" type="DMARCResultType"/>
	    <xs:element name="spf" type="DMARCResultType"/>
	    <xs:element name="reason" type="PolicyOverrideReason"
			minOccurs="0" maxOccurs="unbounded"/>
	  </xs:sequence>
	</xs:complexType>
*/
type PolicyEvaluated struct {
	Disposition     Disposition            `xml:"disposition"`
	Dkim            DMARCResult            `xml:"dkim"`
	Spf             DMARCResult            `xml:"spf"`
	OverrideReasons []PolicyOverrideReason `xml:"reason,omitempty"`
}

/*
A row in the record.

#TODO: Implement the regex for IPAddress as specified in the RFC.

	<xs:complexType name="RowType">
	  <xs:all>
	    <xs:element name="source_ip" type="IPAddress"/>
	    <xs:element name="count" type="xs:integer"/>
	    <xs:element name="policy_evaluated"
			type="PolicyEvaluatedType"
			minOccurs="1"/>
	  </xs:all>
	</xs:complexType>
*/
type Row struct {
	SourceIp        string          `xml:"source_ip"`
	Count           int             `xml:"count"`
	PolicyEvaluated PolicyEvaluated `xml:"policy_evaluated"`
}

/*
The identifying headers from received messages.

Note that while both
envelope_from and header_from are required fields, in practice (from what I've
seen) only header_from is actually used.

	<xs:complexType name="IdentifierType">
	  <xs:all>
	    <xs:element name="envelope_to" type="xs:string" minOccurs="0"/>
	    <xs:element name="envelope_from" type="xs:string" minOccurs="1"/>
	    <xs:element name="header_from" type="xs:string" minOccurs="1"/>
	  </xs:all>
	</xs:complexType>
*/
type Identifier struct {
	EnvelopeTo string `xml:"envelope_to,omitempty"`
	MailFrom   string `xml:"mail_from,omitempty"`
	HeaderFrom string `xml:"header_from"`
}

/*
DKIM verification result. See RFC7001, section 2.6.1.

	<xs:simpleType name="DKIMResultType">
	  <xs:restriction base="xs:string">
	    <xs:enumeration value="none"/>
	    <xs:enumeration value="pass"/>
	    <xs:enumeration value="fail"/>
	    <xs:enumeration value="policy"/>
	    <xs:enumeration value="neutral"/>
	    <xs:enumeration value="temperror"/>
	    <xs:enumeration value="permerror"/>
	  </xs:restriction>
	</xs:simpleType>
*/
type DKIMResult string

/*
The DKIM result, along with information about the "d" and "s" parameters in the
signature.

	<xs:complexType name="DKIMAuthResultType">
	  <xs:all>
	    <xs:element name="domain" type="xs:string" minOccurs="1"/>
	    <xs:element name="selector" type="xs:string" minOccurs="0"/>
	    <xs:element name="result" type="DKIMResultType" minOccurs="1"/>
	    <xs:element name="human_result" type="xs:string" minOccurs="0"/>
	  </xs:all>
	</xs:complexType>
*/
type DKIMAuthResult struct {
	Domain      string     `xml:"domain"`
	Selector    string     `xml:"selector,omitempty"`
	Result      DKIMResult `xml:"result"`
	HumanResult string     `xml:"human_result"`
}

/*
	<xs:simpleType name="SPFDomainScope">
	  <xs:restriction base="xs:string">
	    <xs:enumeration value="helo"/>
	    <xs:enumeration value="mfrom"/>
	  </xs:restriction>
	</xs:simpleType>
*/
type SPFDomainScope string

/*
The result of the SPF check.

	<xs:simpleType name="SPFResultType">
	  <xs:restriction base="xs:string">
	    <xs:enumeration value="none"/>
	    <xs:enumeration value="neutral"/>
	    <xs:enumeration value="pass"/>
	    <xs:enumeration value="fail"/>
	    <xs:enumeration value="softfail"/>
	    <xs:enumeration value="temperror"/>
	    <xs:enumeration value="permerror"/>
	  </xs:restriction>
	</xs:simpleType>
*/
type SPFResult string

/*
The result of the SPF check, with domain and scope information.

	<xs:complexType name="SPFAuthResultType">
	  <xs:all>
	    <xs:element name="domain" type="xs:string" minOccurs="1"/>
	    <xs:element name="scope" type="SPFDomainScope" minOccurs="1"/>
	    <xs:element name="result" type="SPFResultType" minOccurs="1"/>
	  </xs:all>
	</xs:complexType>
*/
type SPFAuthResult struct {
	Domain string         `xml:"domain"`
	Scope  SPFDomainScope `xml:"scope,omitempty"`
	Result SPFResult      `xml:"result"`
}

/*
Combined raw (not influenced by DMARC) DKIM and SPF results.

	<xs:complexType name="AuthResultType">
	  <xs:sequence>
	    <xs:element name="dkim" type="DKIMAuthResultType" minOccurs="0"
			maxOccurs="unbounded"/>
	    <xs:element name="spf" type="SPFAuthResultType" minOccurs="1"
			maxOccurs="unbounded"/>
	  </xs:sequence>
	</xs:complexType>
*/
type AuthResult struct {
	Dkim []DKIMAuthResult `xml:"dkim,omitempty"`
	Spf  []SPFAuthResult  `xml:"spf"`
}

/*
Authentication results from the receiver regarding a certain set of messages.

	<xs:complexType name="RecordType">
	  <xs:sequence>
	    <xs:element name="row" type="RowType"/>
	    <xs:element name="identifiers" type="IdentifierType"/>
	    <xs:element name="auth_results" type="AuthResultType"/>
	  </xs:sequence>
	</xs:complexType>
*/
type Record struct {
	Row        Row        `xml:"row"`
	Identifier Identifier `xml:"identifiers"`
	AuthResult AuthResult `xml:"auth_results"`
}

/*
Parent element tying everything together.

	<xs:element name="feedback">
	  <xs:complexType>
	    <xs:sequence>
	      <xs:element name="version" type="xs:decimal"/>
	      <xs:element name="report_metadata" type="ReportMetadataType"/>
	      <xs:element name="policy_published" type="PolicyPublishedType"/>
	      <xs:element name="record" type="RecordType" maxOccurs="unbounded"/>
	    </xs:sequence>
	  </xs:complexType>
	</xs:element>
*/
type Feedback struct {
	XMLName  xml.Name        `xml:"feedback"`
	Version  float32         `xml:"version,omitempty"`
	Metadata Metadata        `xml:"report_metadata"`
	Policy   PolicyPublished `xml:"policy_published"`
	Records  []Record        `xml:"record"`
}