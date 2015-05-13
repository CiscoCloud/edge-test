package avro

import "github.com/stealthly/go-avro"

type AuthnRequestType struct {
	ForceAuthn                     interface{}
	IsPassive                      interface{}
	ProtocolBinding                interface{}
	AssertionConsumerServiceIndex  interface{}
	AssertionConsumerServiceURL    interface{}
	AttributeConsumingServiceIndex interface{}
	ProviderName                   interface{}
	ID                             string
	Version                        string
	IssueInstant                   string
	Destination                    interface{}
	Consent                        interface{}
	Issuer                         *NameIDType
	Signature                      *SignatureType
	Extensions                     *ExtensionsType
	Subject                        *SubjectType
	NameIDPolicy                   *NameIDPolicyType
	Conditions                     *ConditionsType
	RequestedAuthnContext          *RequestedAuthnContextType
	Scoping                        *ScopingType
}

func NewAuthnRequestType() *AuthnRequestType {
	return &AuthnRequestType{}
}

func (this *AuthnRequestType) Schema() avro.Schema {
	if _AuthnRequestType_schema_err != nil {
		panic(_AuthnRequestType_schema_err)
	}
	return _AuthnRequestType_schema
}

type NameIDType struct {
	NameQualifier   interface{}
	SPNameQualifier interface{}
	Format          interface{}
	SPProvidedID    interface{}
}

func NewNameIDType() *NameIDType {
	return &NameIDType{}
}

func (this *NameIDType) Schema() avro.Schema {
	if _NameIDType_schema_err != nil {
		panic(_NameIDType_schema_err)
	}
	return _NameIDType_schema
}

type SignatureType struct {
	Id             interface{}
	SignedInfo     *SignedInfoType
	SignatureValue *SignatureValueType
	KeyInfo        *KeyInfoType
	Object         []*ObjectType
}

func NewSignatureType() *SignatureType {
	return &SignatureType{
		SignedInfo:     NewSignedInfoType(),
		SignatureValue: NewSignatureValueType(),
		Object:         make([]*ObjectType, 0),
	}
}

func (this *SignatureType) Schema() avro.Schema {
	if _SignatureType_schema_err != nil {
		panic(_SignatureType_schema_err)
	}
	return _SignatureType_schema
}

type SignedInfoType struct {
	Id                     interface{}
	CanonicalizationMethod *CanonicalizationMethodType
	SignatureMethod        *SignatureMethodType
	Reference              []*ReferenceType
}

func NewSignedInfoType() *SignedInfoType {
	return &SignedInfoType{
		CanonicalizationMethod: NewCanonicalizationMethodType(),
		SignatureMethod:        NewSignatureMethodType(),
		Reference:              make([]*ReferenceType, 0),
	}
}

func (this *SignedInfoType) Schema() avro.Schema {
	if _SignedInfoType_schema_err != nil {
		panic(_SignedInfoType_schema_err)
	}
	return _SignedInfoType_schema
}

type CanonicalizationMethodType struct {
	Algorithm string
	Others    map[string]string
}

func NewCanonicalizationMethodType() *CanonicalizationMethodType {
	return &CanonicalizationMethodType{
		Others: make(map[string]string),
	}
}

func (this *CanonicalizationMethodType) Schema() avro.Schema {
	if _CanonicalizationMethodType_schema_err != nil {
		panic(_CanonicalizationMethodType_schema_err)
	}
	return _CanonicalizationMethodType_schema
}

type SignatureMethodType struct {
	Algorithm        string
	HMACOutputLength interface{}
	Others           map[string]string
}

func NewSignatureMethodType() *SignatureMethodType {
	return &SignatureMethodType{
		Others: make(map[string]string),
	}
}

func (this *SignatureMethodType) Schema() avro.Schema {
	if _SignatureMethodType_schema_err != nil {
		panic(_SignatureMethodType_schema_err)
	}
	return _SignatureMethodType_schema
}

type ReferenceType struct {
	Id           interface{}
	URI          interface{}
	Type         interface{}
	Transforms   *TransformsType
	DigestMethod *DigestMethodType
	DigestValue  string
}

func NewReferenceType() *ReferenceType {
	return &ReferenceType{
		DigestMethod: NewDigestMethodType(),
	}
}

func (this *ReferenceType) Schema() avro.Schema {
	if _ReferenceType_schema_err != nil {
		panic(_ReferenceType_schema_err)
	}
	return _ReferenceType_schema
}

type TransformsType struct {
	Transform []*TransformType
}

func NewTransformsType() *TransformsType {
	return &TransformsType{
		Transform: make([]*TransformType, 0),
	}
}

func (this *TransformsType) Schema() avro.Schema {
	if _TransformsType_schema_err != nil {
		panic(_TransformsType_schema_err)
	}
	return _TransformsType_schema
}

type TransformType struct {
	Algorithm string
	Others    map[string]string
	XPath     interface{}
}

func NewTransformType() *TransformType {
	return &TransformType{
		Others: make(map[string]string),
	}
}

func (this *TransformType) Schema() avro.Schema {
	if _TransformType_schema_err != nil {
		panic(_TransformType_schema_err)
	}
	return _TransformType_schema
}

type DigestMethodType struct {
	Algorithm string
	Others    map[string]string
}

func NewDigestMethodType() *DigestMethodType {
	return &DigestMethodType{
		Others: make(map[string]string),
	}
}

func (this *DigestMethodType) Schema() avro.Schema {
	if _DigestMethodType_schema_err != nil {
		panic(_DigestMethodType_schema_err)
	}
	return _DigestMethodType_schema
}

type SignatureValueType struct {
	Id interface{}
}

func NewSignatureValueType() *SignatureValueType {
	return &SignatureValueType{}
}

func (this *SignatureValueType) Schema() avro.Schema {
	if _SignatureValueType_schema_err != nil {
		panic(_SignatureValueType_schema_err)
	}
	return _SignatureValueType_schema
}

type KeyInfoType struct {
	Id              interface{}
	KeyName         interface{}
	KeyValue        *KeyValueType
	RetrievalMethod *RetrievalMethodType
	X509Data        *X509DataType
	PGPData         *PGPDataType
	SPKIData        *SPKIDataType
	MgmtData        interface{}
	Others          map[string]string
}

func NewKeyInfoType() *KeyInfoType {
	return &KeyInfoType{
		Others: make(map[string]string),
	}
}

func (this *KeyInfoType) Schema() avro.Schema {
	if _KeyInfoType_schema_err != nil {
		panic(_KeyInfoType_schema_err)
	}
	return _KeyInfoType_schema
}

type KeyValueType struct {
	DSAKeyValue *DSAKeyValueType
	RSAKeyValue *RSAKeyValueType
	Others      map[string]string
}

func NewKeyValueType() *KeyValueType {
	return &KeyValueType{
		Others: make(map[string]string),
	}
}

func (this *KeyValueType) Schema() avro.Schema {
	if _KeyValueType_schema_err != nil {
		panic(_KeyValueType_schema_err)
	}
	return _KeyValueType_schema
}

type DSAKeyValueType struct {
	P           string
	Q           string
	G           interface{}
	Y           string
	J           interface{}
	Seed        string
	PgenCounter string
}

func NewDSAKeyValueType() *DSAKeyValueType {
	return &DSAKeyValueType{}
}

func (this *DSAKeyValueType) Schema() avro.Schema {
	if _DSAKeyValueType_schema_err != nil {
		panic(_DSAKeyValueType_schema_err)
	}
	return _DSAKeyValueType_schema
}

type RSAKeyValueType struct {
	Modulus  string
	Exponent string
}

func NewRSAKeyValueType() *RSAKeyValueType {
	return &RSAKeyValueType{}
}

func (this *RSAKeyValueType) Schema() avro.Schema {
	if _RSAKeyValueType_schema_err != nil {
		panic(_RSAKeyValueType_schema_err)
	}
	return _RSAKeyValueType_schema
}

type RetrievalMethodType struct {
	URI        interface{}
	Type       interface{}
	Transforms *TransformsType
}

func NewRetrievalMethodType() *RetrievalMethodType {
	return &RetrievalMethodType{}
}

func (this *RetrievalMethodType) Schema() avro.Schema {
	if _RetrievalMethodType_schema_err != nil {
		panic(_RetrievalMethodType_schema_err)
	}
	return _RetrievalMethodType_schema
}

type X509DataType struct {
	X509IssuerSerial *X509IssuerSerialType
	X509SKI          interface{}
	X509SubjectName  interface{}
	X509Certificate  interface{}
	X509CRL          interface{}
	Others           map[string]string
}

func NewX509DataType() *X509DataType {
	return &X509DataType{
		Others: make(map[string]string),
	}
}

func (this *X509DataType) Schema() avro.Schema {
	if _X509DataType_schema_err != nil {
		panic(_X509DataType_schema_err)
	}
	return _X509DataType_schema
}

type X509IssuerSerialType struct {
	X509IssuerName   string
	X509SerialNumber string
}

func NewX509IssuerSerialType() *X509IssuerSerialType {
	return &X509IssuerSerialType{}
}

func (this *X509IssuerSerialType) Schema() avro.Schema {
	if _X509IssuerSerialType_schema_err != nil {
		panic(_X509IssuerSerialType_schema_err)
	}
	return _X509IssuerSerialType_schema
}

type PGPDataType struct {
	PGPKeyID      interface{}
	PGPKeyPacket0 interface{}
	Others        map[string]string
}

func NewPGPDataType() *PGPDataType {
	return &PGPDataType{
		Others: make(map[string]string),
	}
}

func (this *PGPDataType) Schema() avro.Schema {
	if _PGPDataType_schema_err != nil {
		panic(_PGPDataType_schema_err)
	}
	return _PGPDataType_schema
}

type SPKIDataType struct {
	SPKISexp string
	Others   map[string]string
}

func NewSPKIDataType() *SPKIDataType {
	return &SPKIDataType{
		Others: make(map[string]string),
	}
}

func (this *SPKIDataType) Schema() avro.Schema {
	if _SPKIDataType_schema_err != nil {
		panic(_SPKIDataType_schema_err)
	}
	return _SPKIDataType_schema
}

type ObjectType struct {
	Id       interface{}
	MimeType interface{}
	Encoding interface{}
	Others   map[string]string
}

func NewObjectType() *ObjectType {
	return &ObjectType{
		Others: make(map[string]string),
	}
}

func (this *ObjectType) Schema() avro.Schema {
	if _ObjectType_schema_err != nil {
		panic(_ObjectType_schema_err)
	}
	return _ObjectType_schema
}

type ExtensionsType struct {
	Others map[string]string
}

func NewExtensionsType() *ExtensionsType {
	return &ExtensionsType{
		Others: make(map[string]string),
	}
}

func (this *ExtensionsType) Schema() avro.Schema {
	if _ExtensionsType_schema_err != nil {
		panic(_ExtensionsType_schema_err)
	}
	return _ExtensionsType_schema
}

type SubjectType struct {
	BaseID               *BaseIDAbstractType
	NameID               *NameIDType
	EncryptedID          *EncryptedElementType
	SubjectConfirmation0 []*SubjectConfirmationType
}

func NewSubjectType() *SubjectType {
	return &SubjectType{
		SubjectConfirmation0: make([]*SubjectConfirmationType, 0),
	}
}

func (this *SubjectType) Schema() avro.Schema {
	if _SubjectType_schema_err != nil {
		panic(_SubjectType_schema_err)
	}
	return _SubjectType_schema
}

type BaseIDAbstractType struct {
	NameQualifier   interface{}
	SPNameQualifier interface{}
}

func NewBaseIDAbstractType() *BaseIDAbstractType {
	return &BaseIDAbstractType{}
}

func (this *BaseIDAbstractType) Schema() avro.Schema {
	if _BaseIDAbstractType_schema_err != nil {
		panic(_BaseIDAbstractType_schema_err)
	}
	return _BaseIDAbstractType_schema
}

type EncryptedElementType struct {
	EncryptedData *EncryptedDataType
	EncryptedKey  []*EncryptedKeyType
}

func NewEncryptedElementType() *EncryptedElementType {
	return &EncryptedElementType{
		EncryptedData: NewEncryptedDataType(),
		EncryptedKey:  make([]*EncryptedKeyType, 0),
	}
}

func (this *EncryptedElementType) Schema() avro.Schema {
	if _EncryptedElementType_schema_err != nil {
		panic(_EncryptedElementType_schema_err)
	}
	return _EncryptedElementType_schema
}

type EncryptedDataType struct {
	Id                   interface{}
	Type                 interface{}
	MimeType             interface{}
	Encoding             interface{}
	EncryptionMethod     *EncryptionMethodType
	KeyInfo              *KeyInfoType
	CipherData           *CipherDataType
	EncryptionProperties *EncryptionPropertiesType
}

func NewEncryptedDataType() *EncryptedDataType {
	return &EncryptedDataType{
		CipherData: NewCipherDataType(),
	}
}

func (this *EncryptedDataType) Schema() avro.Schema {
	if _EncryptedDataType_schema_err != nil {
		panic(_EncryptedDataType_schema_err)
	}
	return _EncryptedDataType_schema
}

type EncryptionMethodType struct {
	Algorithm  string
	KeySize    interface{}
	OAEPparams interface{}
	Others     map[string]string
}

func NewEncryptionMethodType() *EncryptionMethodType {
	return &EncryptionMethodType{
		Others: make(map[string]string),
	}
}

func (this *EncryptionMethodType) Schema() avro.Schema {
	if _EncryptionMethodType_schema_err != nil {
		panic(_EncryptionMethodType_schema_err)
	}
	return _EncryptionMethodType_schema
}

type CipherDataType struct {
	CipherValue     interface{}
	CipherReference *CipherReferenceType
}

func NewCipherDataType() *CipherDataType {
	return &CipherDataType{}
}

func (this *CipherDataType) Schema() avro.Schema {
	if _CipherDataType_schema_err != nil {
		panic(_CipherDataType_schema_err)
	}
	return _CipherDataType_schema
}

type CipherReferenceType struct {
	URI        string
	Transforms *TransformsType
}

func NewCipherReferenceType() *CipherReferenceType {
	return &CipherReferenceType{}
}

func (this *CipherReferenceType) Schema() avro.Schema {
	if _CipherReferenceType_schema_err != nil {
		panic(_CipherReferenceType_schema_err)
	}
	return _CipherReferenceType_schema
}

type EncryptionPropertiesType struct {
	Id                 interface{}
	EncryptionProperty []*EncryptionPropertyType
}

func NewEncryptionPropertiesType() *EncryptionPropertiesType {
	return &EncryptionPropertiesType{
		EncryptionProperty: make([]*EncryptionPropertyType, 0),
	}
}

func (this *EncryptionPropertiesType) Schema() avro.Schema {
	if _EncryptionPropertiesType_schema_err != nil {
		panic(_EncryptionPropertiesType_schema_err)
	}
	return _EncryptionPropertiesType_schema
}

type EncryptionPropertyType struct {
	Target interface{}
	Id     interface{}
	Others map[string]string
}

func NewEncryptionPropertyType() *EncryptionPropertyType {
	return &EncryptionPropertyType{
		Others: make(map[string]string),
	}
}

func (this *EncryptionPropertyType) Schema() avro.Schema {
	if _EncryptionPropertyType_schema_err != nil {
		panic(_EncryptionPropertyType_schema_err)
	}
	return _EncryptionPropertyType_schema
}

type EncryptedKeyType struct {
	Recipient            interface{}
	Id                   interface{}
	Type                 interface{}
	MimeType             interface{}
	Encoding             interface{}
	EncryptionMethod     *EncryptionMethodType
	KeyInfo              *KeyInfoType
	CipherData           *CipherDataType
	EncryptionProperties *EncryptionPropertiesType
	ReferenceList        *Type0
	CarriedKeyName       interface{}
}

func NewEncryptedKeyType() *EncryptedKeyType {
	return &EncryptedKeyType{}
}

func (this *EncryptedKeyType) Schema() avro.Schema {
	if _EncryptedKeyType_schema_err != nil {
		panic(_EncryptedKeyType_schema_err)
	}
	return _EncryptedKeyType_schema
}

type Type0 struct {
	DataReference *ReferenceType
	KeyReference  *ReferenceType
}

func NewType0() *Type0 {
	return &Type0{}
}

func (this *Type0) Schema() avro.Schema {
	if _Type0_schema_err != nil {
		panic(_Type0_schema_err)
	}
	return _Type0_schema
}

type SubjectConfirmationType struct {
	Method                  string
	BaseID                  *BaseIDAbstractType
	NameID                  *NameIDType
	EncryptedID             *EncryptedElementType
	SubjectConfirmationData *SubjectConfirmationDataType
}

func NewSubjectConfirmationType() *SubjectConfirmationType {
	return &SubjectConfirmationType{}
}

func (this *SubjectConfirmationType) Schema() avro.Schema {
	if _SubjectConfirmationType_schema_err != nil {
		panic(_SubjectConfirmationType_schema_err)
	}
	return _SubjectConfirmationType_schema
}

type SubjectConfirmationDataType struct {
	NotBefore    interface{}
	NotOnOrAfter interface{}
	Recipient    interface{}
	InResponseTo interface{}
	Address      interface{}
	Others       map[string]string
}

func NewSubjectConfirmationDataType() *SubjectConfirmationDataType {
	return &SubjectConfirmationDataType{
		Others: make(map[string]string),
	}
}

func (this *SubjectConfirmationDataType) Schema() avro.Schema {
	if _SubjectConfirmationDataType_schema_err != nil {
		panic(_SubjectConfirmationDataType_schema_err)
	}
	return _SubjectConfirmationDataType_schema
}

type NameIDPolicyType struct {
	Format          interface{}
	SPNameQualifier interface{}
	AllowCreate     interface{}
}

func NewNameIDPolicyType() *NameIDPolicyType {
	return &NameIDPolicyType{}
}

func (this *NameIDPolicyType) Schema() avro.Schema {
	if _NameIDPolicyType_schema_err != nil {
		panic(_NameIDPolicyType_schema_err)
	}
	return _NameIDPolicyType_schema
}

type ConditionsType struct {
	NotBefore           interface{}
	NotOnOrAfter        interface{}
	Condition           *ConditionAbstractType
	AudienceRestriction *AudienceRestrictionType
	OneTimeUse          *OneTimeUseType
	ProxyRestriction    *ProxyRestrictionType
}

func NewConditionsType() *ConditionsType {
	return &ConditionsType{}
}

func (this *ConditionsType) Schema() avro.Schema {
	if _ConditionsType_schema_err != nil {
		panic(_ConditionsType_schema_err)
	}
	return _ConditionsType_schema
}

type ConditionAbstractType struct {
}

func NewConditionAbstractType() *ConditionAbstractType {
	return &ConditionAbstractType{}
}

func (this *ConditionAbstractType) Schema() avro.Schema {
	if _ConditionAbstractType_schema_err != nil {
		panic(_ConditionAbstractType_schema_err)
	}
	return _ConditionAbstractType_schema
}

type AudienceRestrictionType struct {
	Audience []string
}

func NewAudienceRestrictionType() *AudienceRestrictionType {
	return &AudienceRestrictionType{
		Audience: make([]string, 0),
	}
}

func (this *AudienceRestrictionType) Schema() avro.Schema {
	if _AudienceRestrictionType_schema_err != nil {
		panic(_AudienceRestrictionType_schema_err)
	}
	return _AudienceRestrictionType_schema
}

type OneTimeUseType struct {
}

func NewOneTimeUseType() *OneTimeUseType {
	return &OneTimeUseType{}
}

func (this *OneTimeUseType) Schema() avro.Schema {
	if _OneTimeUseType_schema_err != nil {
		panic(_OneTimeUseType_schema_err)
	}
	return _OneTimeUseType_schema
}

type ProxyRestrictionType struct {
	Count    interface{}
	Audience []string
}

func NewProxyRestrictionType() *ProxyRestrictionType {
	return &ProxyRestrictionType{
		Audience: make([]string, 0),
	}
}

func (this *ProxyRestrictionType) Schema() avro.Schema {
	if _ProxyRestrictionType_schema_err != nil {
		panic(_ProxyRestrictionType_schema_err)
	}
	return _ProxyRestrictionType_schema
}

type RequestedAuthnContextType struct {
	Comparison           interface{}
	AuthnContextClassRef []string
	AuthnContextDeclRef  []string
}

func NewRequestedAuthnContextType() *RequestedAuthnContextType {
	return &RequestedAuthnContextType{
		AuthnContextClassRef: make([]string, 0),
		AuthnContextDeclRef:  make([]string, 0),
	}
}

func (this *RequestedAuthnContextType) Schema() avro.Schema {
	if _RequestedAuthnContextType_schema_err != nil {
		panic(_RequestedAuthnContextType_schema_err)
	}
	return _RequestedAuthnContextType_schema
}

type ScopingType struct {
	ProxyCount  interface{}
	IDPList     *IDPListType
	RequesterID []string
}

func NewScopingType() *ScopingType {
	return &ScopingType{
		RequesterID: make([]string, 0),
	}
}

func (this *ScopingType) Schema() avro.Schema {
	if _ScopingType_schema_err != nil {
		panic(_ScopingType_schema_err)
	}
	return _ScopingType_schema
}

type IDPListType struct {
	IDPEntry    []*IDPEntryType
	GetComplete interface{}
}

func NewIDPListType() *IDPListType {
	return &IDPListType{
		IDPEntry: make([]*IDPEntryType, 0),
	}
}

func (this *IDPListType) Schema() avro.Schema {
	if _IDPListType_schema_err != nil {
		panic(_IDPListType_schema_err)
	}
	return _IDPListType_schema
}

type IDPEntryType struct {
	ProviderID string
	Name       interface{}
	Loc        interface{}
}

func NewIDPEntryType() *IDPEntryType {
	return &IDPEntryType{}
}

func (this *IDPEntryType) Schema() avro.Schema {
	if _IDPEntryType_schema_err != nil {
		panic(_IDPEntryType_schema_err)
	}
	return _IDPEntryType_schema
}

type ArtifactResponseType struct {
	ID           string
	InResponseTo interface{}
	Version      string
	IssueInstant string
	Destination  interface{}
	Consent      interface{}
	Issuer       *NameIDType
	Signature    *SignatureType
	Extensions   *ExtensionsType
	Status       *StatusType
	Others       map[string]string
}

func NewArtifactResponseType() *ArtifactResponseType {
	return &ArtifactResponseType{
		Status: NewStatusType(),
		Others: make(map[string]string),
	}
}

func (this *ArtifactResponseType) Schema() avro.Schema {
	if _ArtifactResponseType_schema_err != nil {
		panic(_ArtifactResponseType_schema_err)
	}
	return _ArtifactResponseType_schema
}

type StatusType struct {
	StatusCode    *StatusCodeType
	StatusMessage interface{}
	StatusDetail  *StatusDetailType
}

func NewStatusType() *StatusType {
	return &StatusType{
		StatusCode: NewStatusCodeType(),
	}
}

func (this *StatusType) Schema() avro.Schema {
	if _StatusType_schema_err != nil {
		panic(_StatusType_schema_err)
	}
	return _StatusType_schema
}

type StatusCodeType struct {
	Value      string
	StatusCode *StatusCodeType
}

func NewStatusCodeType() *StatusCodeType {
	return &StatusCodeType{}
}

func (this *StatusCodeType) Schema() avro.Schema {
	if _StatusCodeType_schema_err != nil {
		panic(_StatusCodeType_schema_err)
	}
	return _StatusCodeType_schema
}

type StatusDetailType struct {
	Others map[string]string
}

func NewStatusDetailType() *StatusDetailType {
	return &StatusDetailType{
		Others: make(map[string]string),
	}
}

func (this *StatusDetailType) Schema() avro.Schema {
	if _StatusDetailType_schema_err != nil {
		panic(_StatusDetailType_schema_err)
	}
	return _StatusDetailType_schema
}

// Generated by codegen. Please do not modify.
var _AuthnRequestType_schema, _AuthnRequestType_schema_err = avro.ParseSchema(`{
    "type": "record",
    "namespace": "avro",
    "name": "AuthnRequestType",
    "fields": [
        {
            "name": "ForceAuthn",
            "type": [
                "boolean",
                "null"
            ]
        },
        {
            "name": "IsPassive",
            "type": [
                "boolean",
                "null"
            ]
        },
        {
            "name": "ProtocolBinding",
            "type": [
                "string",
                "null"
            ]
        },
        {
            "name": "AssertionConsumerServiceIndex",
            "type": [
                "string",
                "null"
            ]
        },
        {
            "name": "AssertionConsumerServiceURL",
            "type": [
                "string",
                "null"
            ]
        },
        {
            "name": "AttributeConsumingServiceIndex",
            "type": [
                "string",
                "null"
            ]
        },
        {
            "name": "ProviderName",
            "type": [
                "string",
                "null"
            ]
        },
        {
            "name": "ID",
            "type": "string"
        },
        {
            "name": "Version",
            "type": "string"
        },
        {
            "name": "IssueInstant",
            "type": "string"
        },
        {
            "name": "Destination",
            "type": [
                "string",
                "null"
            ]
        },
        {
            "name": "Consent",
            "type": [
                "string",
                "null"
            ]
        },
        {
            "name": "Issuer",
            "type": [
                {
                    "type": "record",
                    "name": "NameIDType",
                    "fields": [
                        {
                            "name": "NameQualifier",
                            "type": [
                                "string",
                                "null"
                            ]
                        },
                        {
                            "name": "SPNameQualifier",
                            "type": [
                                "string",
                                "null"
                            ]
                        },
                        {
                            "name": "Format",
                            "type": [
                                "string",
                                "null"
                            ]
                        },
                        {
                            "name": "SPProvidedID",
                            "type": [
                                "string",
                                "null"
                            ]
                        }
                    ]
                },
                "null"
            ]
        },
        {
            "name": "Signature",
            "type": [
                {
                    "type": "record",
                    "name": "SignatureType",
                    "fields": [
                        {
                            "name": "Id",
                            "type": [
                                "string",
                                "null"
                            ]
                        },
                        {
                            "name": "SignedInfo",
                            "type": {
                                "type": "record",
                                "name": "SignedInfoType",
                                "fields": [
                                    {
                                        "name": "Id",
                                        "type": [
                                            "string",
                                            "null"
                                        ]
                                    },
                                    {
                                        "name": "CanonicalizationMethod",
                                        "type": {
                                            "type": "record",
                                            "name": "CanonicalizationMethodType",
                                            "fields": [
                                                {
                                                    "name": "Algorithm",
                                                    "type": "string"
                                                },
                                                {
                                                    "name": "others",
                                                    "type": {
                                                        "type": "map",
                                                        "values": "string"
                                                    }
                                                }
                                            ]
                                        }
                                    },
                                    {
                                        "name": "SignatureMethod",
                                        "type": {
                                            "type": "record",
                                            "name": "SignatureMethodType",
                                            "fields": [
                                                {
                                                    "name": "Algorithm",
                                                    "type": "string"
                                                },
                                                {
                                                    "name": "HMACOutputLength",
                                                    "type": [
                                                        "string",
                                                        "null"
                                                    ]
                                                },
                                                {
                                                    "name": "others",
                                                    "type": {
                                                        "type": "map",
                                                        "values": "string"
                                                    }
                                                }
                                            ]
                                        }
                                    },
                                    {
                                        "name": "Reference",
                                        "type": {
                                            "type": "array",
                                            "items": {
                                                "type": "record",
                                                "name": "ReferenceType",
                                                "fields": [
                                                    {
                                                        "name": "Id",
                                                        "type": [
                                                            "string",
                                                            "null"
                                                        ]
                                                    },
                                                    {
                                                        "name": "URI",
                                                        "type": [
                                                            "string",
                                                            "null"
                                                        ]
                                                    },
                                                    {
                                                        "name": "Type",
                                                        "type": [
                                                            "string",
                                                            "null"
                                                        ]
                                                    },
                                                    {
                                                        "name": "Transforms",
                                                        "type": [
                                                            {
                                                                "type": "record",
                                                                "name": "TransformsType",
                                                                "fields": [
                                                                    {
                                                                        "name": "Transform",
                                                                        "type": {
                                                                            "type": "array",
                                                                            "items": {
                                                                                "type": "record",
                                                                                "name": "TransformType",
                                                                                "fields": [
                                                                                    {
                                                                                        "name": "Algorithm",
                                                                                        "type": "string"
                                                                                    },
                                                                                    {
                                                                                        "name": "others",
                                                                                        "type": {
                                                                                            "type": "map",
                                                                                            "values": "string"
                                                                                        }
                                                                                    },
                                                                                    {
                                                                                        "name": "XPath",
                                                                                        "type": [
                                                                                            "string",
                                                                                            "null"
                                                                                        ]
                                                                                    }
                                                                                ]
                                                                            }
                                                                        }
                                                                    }
                                                                ]
                                                            },
                                                            "null"
                                                        ]
                                                    },
                                                    {
                                                        "name": "DigestMethod",
                                                        "type": {
                                                            "type": "record",
                                                            "name": "DigestMethodType",
                                                            "fields": [
                                                                {
                                                                    "name": "Algorithm",
                                                                    "type": "string"
                                                                },
                                                                {
                                                                    "name": "others",
                                                                    "type": {
                                                                        "type": "map",
                                                                        "values": "string"
                                                                    }
                                                                }
                                                            ]
                                                        }
                                                    },
                                                    {
                                                        "name": "DigestValue",
                                                        "type": "string"
                                                    }
                                                ]
                                            }
                                        }
                                    }
                                ]
                            }
                        },
                        {
                            "name": "SignatureValue",
                            "type": {
                                "type": "record",
                                "name": "SignatureValueType",
                                "fields": [
                                    {
                                        "name": "Id",
                                        "type": [
                                            "string",
                                            "null"
                                        ]
                                    }
                                ]
                            }
                        },
                        {
                            "name": "KeyInfo",
                            "type": [
                                {
                                    "type": "record",
                                    "name": "KeyInfoType",
                                    "fields": [
                                        {
                                            "name": "Id",
                                            "type": [
                                                "string",
                                                "null"
                                            ]
                                        },
                                        {
                                            "name": "KeyName",
                                            "type": [
                                                "string",
                                                "null"
                                            ]
                                        },
                                        {
                                            "name": "KeyValue",
                                            "type": [
                                                {
                                                    "type": "record",
                                                    "name": "KeyValueType",
                                                    "fields": [
                                                        {
                                                            "name": "DSAKeyValue",
                                                            "type": [
                                                                {
                                                                    "type": "record",
                                                                    "name": "DSAKeyValueType",
                                                                    "fields": [
                                                                        {
                                                                            "name": "P",
                                                                            "type": "string"
                                                                        },
                                                                        {
                                                                            "name": "Q",
                                                                            "type": "string"
                                                                        },
                                                                        {
                                                                            "name": "G",
                                                                            "type": [
                                                                                "string",
                                                                                "null"
                                                                            ]
                                                                        },
                                                                        {
                                                                            "name": "Y",
                                                                            "type": "string"
                                                                        },
                                                                        {
                                                                            "name": "J",
                                                                            "type": [
                                                                                "string",
                                                                                "null"
                                                                            ]
                                                                        },
                                                                        {
                                                                            "name": "Seed",
                                                                            "type": "string"
                                                                        },
                                                                        {
                                                                            "name": "PgenCounter",
                                                                            "type": "string"
                                                                        }
                                                                    ]
                                                                },
                                                                "null"
                                                            ]
                                                        },
                                                        {
                                                            "name": "RSAKeyValue",
                                                            "type": [
                                                                {
                                                                    "type": "record",
                                                                    "name": "RSAKeyValueType",
                                                                    "fields": [
                                                                        {
                                                                            "name": "Modulus",
                                                                            "type": "string"
                                                                        },
                                                                        {
                                                                            "name": "Exponent",
                                                                            "type": "string"
                                                                        }
                                                                    ]
                                                                },
                                                                "null"
                                                            ]
                                                        },
                                                        {
                                                            "name": "others",
                                                            "type": {
                                                                "type": "map",
                                                                "values": "string"
                                                            }
                                                        }
                                                    ]
                                                },
                                                "null"
                                            ]
                                        },
                                        {
                                            "name": "RetrievalMethod",
                                            "type": [
                                                {
                                                    "type": "record",
                                                    "name": "RetrievalMethodType",
                                                    "fields": [
                                                        {
                                                            "name": "URI",
                                                            "type": [
                                                                "string",
                                                                "null"
                                                            ]
                                                        },
                                                        {
                                                            "name": "Type",
                                                            "type": [
                                                                "string",
                                                                "null"
                                                            ]
                                                        },
                                                        {
                                                            "name": "Transforms",
                                                            "type": [
                                                                "TransformsType",
                                                                "null"
                                                            ]
                                                        }
                                                    ]
                                                },
                                                "null"
                                            ]
                                        },
                                        {
                                            "name": "X509Data",
                                            "type": [
                                                {
                                                    "type": "record",
                                                    "name": "X509DataType",
                                                    "fields": [
                                                        {
                                                            "name": "X509IssuerSerial",
                                                            "type": [
                                                                {
                                                                    "type": "record",
                                                                    "name": "X509IssuerSerialType",
                                                                    "fields": [
                                                                        {
                                                                            "name": "X509IssuerName",
                                                                            "type": "string"
                                                                        },
                                                                        {
                                                                            "name": "X509SerialNumber",
                                                                            "type": "string"
                                                                        }
                                                                    ]
                                                                },
                                                                "null"
                                                            ]
                                                        },
                                                        {
                                                            "name": "X509SKI",
                                                            "type": [
                                                                "string",
                                                                "null"
                                                            ]
                                                        },
                                                        {
                                                            "name": "X509SubjectName",
                                                            "type": [
                                                                "string",
                                                                "null"
                                                            ]
                                                        },
                                                        {
                                                            "name": "X509Certificate",
                                                            "type": [
                                                                "string",
                                                                "null"
                                                            ]
                                                        },
                                                        {
                                                            "name": "X509CRL",
                                                            "type": [
                                                                "string",
                                                                "null"
                                                            ]
                                                        },
                                                        {
                                                            "name": "others",
                                                            "type": {
                                                                "type": "map",
                                                                "values": "string"
                                                            }
                                                        }
                                                    ]
                                                },
                                                "null"
                                            ]
                                        },
                                        {
                                            "name": "PGPData",
                                            "type": [
                                                {
                                                    "type": "record",
                                                    "name": "PGPDataType",
                                                    "fields": [
                                                        {
                                                            "name": "PGPKeyID",
                                                            "type": [
                                                                "string",
                                                                "null"
                                                            ]
                                                        },
                                                        {
                                                            "name": "PGPKeyPacket0",
                                                            "type": [
                                                                "string",
                                                                "null"
                                                            ]
                                                        },
                                                        {
                                                            "name": "others",
                                                            "type": {
                                                                "type": "map",
                                                                "values": "string"
                                                            }
                                                        }
                                                    ]
                                                },
                                                "null"
                                            ]
                                        },
                                        {
                                            "name": "SPKIData",
                                            "type": [
                                                {
                                                    "type": "record",
                                                    "name": "SPKIDataType",
                                                    "fields": [
                                                        {
                                                            "name": "SPKISexp",
                                                            "type": "string"
                                                        },
                                                        {
                                                            "name": "others",
                                                            "type": {
                                                                "type": "map",
                                                                "values": "string"
                                                            }
                                                        }
                                                    ]
                                                },
                                                "null"
                                            ]
                                        },
                                        {
                                            "name": "MgmtData",
                                            "type": [
                                                "string",
                                                "null"
                                            ]
                                        },
                                        {
                                            "name": "others",
                                            "type": {
                                                "type": "map",
                                                "values": "string"
                                            }
                                        }
                                    ]
                                },
                                "null"
                            ]
                        },
                        {
                            "name": "Object",
                            "type": {
                                "type": "array",
                                "items": {
                                    "type": "record",
                                    "name": "ObjectType",
                                    "fields": [
                                        {
                                            "name": "Id",
                                            "type": [
                                                "string",
                                                "null"
                                            ]
                                        },
                                        {
                                            "name": "MimeType",
                                            "type": [
                                                "string",
                                                "null"
                                            ]
                                        },
                                        {
                                            "name": "Encoding",
                                            "type": [
                                                "string",
                                                "null"
                                            ]
                                        },
                                        {
                                            "name": "others",
                                            "type": {
                                                "type": "map",
                                                "values": "string"
                                            }
                                        }
                                    ]
                                }
                            }
                        }
                    ]
                },
                "null"
            ]
        },
        {
            "name": "Extensions",
            "type": [
                {
                    "type": "record",
                    "name": "ExtensionsType",
                    "fields": [
                        {
                            "name": "others",
                            "type": {
                                "type": "map",
                                "values": "string"
                            }
                        }
                    ]
                },
                "null"
            ]
        },
        {
            "name": "Subject",
            "type": [
                {
                    "type": "record",
                    "name": "SubjectType",
                    "fields": [
                        {
                            "name": "BaseID",
                            "type": [
                                {
                                    "type": "record",
                                    "name": "BaseIDAbstractType",
                                    "fields": [
                                        {
                                            "name": "NameQualifier",
                                            "type": [
                                                "string",
                                                "null"
                                            ]
                                        },
                                        {
                                            "name": "SPNameQualifier",
                                            "type": [
                                                "string",
                                                "null"
                                            ]
                                        }
                                    ]
                                },
                                "null"
                            ]
                        },
                        {
                            "name": "NameID",
                            "type": [
                                "NameIDType",
                                "null"
                            ]
                        },
                        {
                            "name": "EncryptedID",
                            "type": [
                                {
                                    "type": "record",
                                    "name": "EncryptedElementType",
                                    "fields": [
                                        {
                                            "name": "EncryptedData",
                                            "type": {
                                                "type": "record",
                                                "name": "EncryptedDataType",
                                                "fields": [
                                                    {
                                                        "name": "Id",
                                                        "type": [
                                                            "string",
                                                            "null"
                                                        ]
                                                    },
                                                    {
                                                        "name": "Type",
                                                        "type": [
                                                            "string",
                                                            "null"
                                                        ]
                                                    },
                                                    {
                                                        "name": "MimeType",
                                                        "type": [
                                                            "string",
                                                            "null"
                                                        ]
                                                    },
                                                    {
                                                        "name": "Encoding",
                                                        "type": [
                                                            "string",
                                                            "null"
                                                        ]
                                                    },
                                                    {
                                                        "name": "EncryptionMethod",
                                                        "type": [
                                                            {
                                                                "type": "record",
                                                                "name": "EncryptionMethodType",
                                                                "fields": [
                                                                    {
                                                                        "name": "Algorithm",
                                                                        "type": "string"
                                                                    },
                                                                    {
                                                                        "name": "KeySize",
                                                                        "type": [
                                                                            "string",
                                                                            "null"
                                                                        ]
                                                                    },
                                                                    {
                                                                        "name": "OAEPparams",
                                                                        "type": [
                                                                            "string",
                                                                            "null"
                                                                        ]
                                                                    },
                                                                    {
                                                                        "name": "others",
                                                                        "type": {
                                                                            "type": "map",
                                                                            "values": "string"
                                                                        }
                                                                    }
                                                                ]
                                                            },
                                                            "null"
                                                        ]
                                                    },
                                                    {
                                                        "name": "KeyInfo",
                                                        "type": [
                                                            "KeyInfoType",
                                                            "null"
                                                        ]
                                                    },
                                                    {
                                                        "name": "CipherData",
                                                        "type": {
                                                            "type": "record",
                                                            "name": "CipherDataType",
                                                            "fields": [
                                                                {
                                                                    "name": "CipherValue",
                                                                    "type": [
                                                                        "string",
                                                                        "null"
                                                                    ]
                                                                },
                                                                {
                                                                    "name": "CipherReference",
                                                                    "type": [
                                                                        {
                                                                            "type": "record",
                                                                            "name": "CipherReferenceType",
                                                                            "fields": [
                                                                                {
                                                                                    "name": "URI",
                                                                                    "type": "string"
                                                                                },
                                                                                {
                                                                                    "name": "Transforms",
                                                                                    "type": [
                                                                                        "TransformsType",
                                                                                        "null"
                                                                                    ]
                                                                                }
                                                                            ]
                                                                        },
                                                                        "null"
                                                                    ]
                                                                }
                                                            ]
                                                        }
                                                    },
                                                    {
                                                        "name": "EncryptionProperties",
                                                        "type": [
                                                            {
                                                                "type": "record",
                                                                "name": "EncryptionPropertiesType",
                                                                "fields": [
                                                                    {
                                                                        "name": "Id",
                                                                        "type": [
                                                                            "string",
                                                                            "null"
                                                                        ]
                                                                    },
                                                                    {
                                                                        "name": "EncryptionProperty",
                                                                        "type": {
                                                                            "type": "array",
                                                                            "items": {
                                                                                "type": "record",
                                                                                "name": "EncryptionPropertyType",
                                                                                "fields": [
                                                                                    {
                                                                                        "name": "Target",
                                                                                        "type": [
                                                                                            "string",
                                                                                            "null"
                                                                                        ]
                                                                                    },
                                                                                    {
                                                                                        "name": "Id",
                                                                                        "type": [
                                                                                            "string",
                                                                                            "null"
                                                                                        ]
                                                                                    },
                                                                                    {
                                                                                        "name": "others",
                                                                                        "type": {
                                                                                            "type": "map",
                                                                                            "values": "string"
                                                                                        }
                                                                                    }
                                                                                ]
                                                                            }
                                                                        }
                                                                    }
                                                                ]
                                                            },
                                                            "null"
                                                        ]
                                                    }
                                                ]
                                            }
                                        },
                                        {
                                            "name": "EncryptedKey",
                                            "type": {
                                                "type": "array",
                                                "items": {
                                                    "type": "record",
                                                    "name": "EncryptedKeyType",
                                                    "fields": [
                                                        {
                                                            "name": "Recipient",
                                                            "type": [
                                                                "string",
                                                                "null"
                                                            ]
                                                        },
                                                        {
                                                            "name": "Id",
                                                            "type": [
                                                                "string",
                                                                "null"
                                                            ]
                                                        },
                                                        {
                                                            "name": "Type",
                                                            "type": [
                                                                "string",
                                                                "null"
                                                            ]
                                                        },
                                                        {
                                                            "name": "MimeType",
                                                            "type": [
                                                                "string",
                                                                "null"
                                                            ]
                                                        },
                                                        {
                                                            "name": "Encoding",
                                                            "type": [
                                                                "string",
                                                                "null"
                                                            ]
                                                        },
                                                        {
                                                            "name": "EncryptionMethod",
                                                            "type": [
                                                                "EncryptionMethodType",
                                                                "null"
                                                            ]
                                                        },
                                                        {
                                                            "name": "KeyInfo",
                                                            "type": [
                                                                "KeyInfoType",
                                                                "null"
                                                            ]
                                                        },
                                                        {
                                                            "name": "CipherData",
                                                            "type": "CipherDataType"
                                                        },
                                                        {
                                                            "name": "EncryptionProperties",
                                                            "type": [
                                                                "EncryptionPropertiesType",
                                                                "null"
                                                            ]
                                                        },
                                                        {
                                                            "name": "ReferenceList",
                                                            "type": [
                                                                {
                                                                    "type": "record",
                                                                    "name": "type0",
                                                                    "fields": [
                                                                        {
                                                                            "name": "DataReference",
                                                                            "type": [
                                                                                "ReferenceType",
                                                                                "null"
                                                                            ]
                                                                        },
                                                                        {
                                                                            "name": "KeyReference",
                                                                            "type": [
                                                                                "ReferenceType",
                                                                                "null"
                                                                            ]
                                                                        }
                                                                    ]
                                                                },
                                                                "null"
                                                            ]
                                                        },
                                                        {
                                                            "name": "CarriedKeyName",
                                                            "type": [
                                                                "string",
                                                                "null"
                                                            ]
                                                        }
                                                    ]
                                                }
                                            }
                                        }
                                    ]
                                },
                                "null"
                            ]
                        },
                        {
                            "name": "SubjectConfirmation0",
                            "type": {
                                "type": "array",
                                "items": {
                                    "type": "record",
                                    "name": "SubjectConfirmationType",
                                    "fields": [
                                        {
                                            "name": "Method",
                                            "type": "string"
                                        },
                                        {
                                            "name": "BaseID",
                                            "type": [
                                                "BaseIDAbstractType",
                                                "null"
                                            ]
                                        },
                                        {
                                            "name": "NameID",
                                            "type": [
                                                "NameIDType",
                                                "null"
                                            ]
                                        },
                                        {
                                            "name": "EncryptedID",
                                            "type": [
                                                "EncryptedElementType",
                                                "null"
                                            ]
                                        },
                                        {
                                            "name": "SubjectConfirmationData",
                                            "type": [
                                                {
                                                    "type": "record",
                                                    "name": "SubjectConfirmationDataType",
                                                    "fields": [
                                                        {
                                                            "name": "NotBefore",
                                                            "type": [
                                                                "string",
                                                                "null"
                                                            ]
                                                        },
                                                        {
                                                            "name": "NotOnOrAfter",
                                                            "type": [
                                                                "string",
                                                                "null"
                                                            ]
                                                        },
                                                        {
                                                            "name": "Recipient",
                                                            "type": [
                                                                "string",
                                                                "null"
                                                            ]
                                                        },
                                                        {
                                                            "name": "InResponseTo",
                                                            "type": [
                                                                "string",
                                                                "null"
                                                            ]
                                                        },
                                                        {
                                                            "name": "Address",
                                                            "type": [
                                                                "string",
                                                                "null"
                                                            ]
                                                        },
                                                        {
                                                            "name": "others",
                                                            "type": {
                                                                "type": "map",
                                                                "values": "string"
                                                            }
                                                        }
                                                    ]
                                                },
                                                "null"
                                            ]
                                        }
                                    ]
                                }
                            }
                        }
                    ]
                },
                "null"
            ]
        },
        {
            "name": "NameIDPolicy",
            "type": [
                {
                    "type": "record",
                    "name": "NameIDPolicyType",
                    "fields": [
                        {
                            "name": "Format",
                            "type": [
                                "string",
                                "null"
                            ]
                        },
                        {
                            "name": "SPNameQualifier",
                            "type": [
                                "string",
                                "null"
                            ]
                        },
                        {
                            "name": "AllowCreate",
                            "type": [
                                "boolean",
                                "null"
                            ]
                        }
                    ]
                },
                "null"
            ]
        },
        {
            "name": "Conditions",
            "type": [
                {
                    "type": "record",
                    "name": "ConditionsType",
                    "fields": [
                        {
                            "name": "NotBefore",
                            "type": [
                                "string",
                                "null"
                            ]
                        },
                        {
                            "name": "NotOnOrAfter",
                            "type": [
                                "string",
                                "null"
                            ]
                        },
                        {
                            "name": "Condition",
                            "type": [
                                {
                                    "type": "record",
                                    "name": "ConditionAbstractType"
                                },
                                "null"
                            ]
                        },
                        {
                            "name": "AudienceRestriction",
                            "type": [
                                {
                                    "type": "record",
                                    "name": "AudienceRestrictionType",
                                    "fields": [
                                        {
                                            "name": "Audience",
                                            "type": {
                                                "type": "array",
                                                "items": "string"
                                            }
                                        }
                                    ]
                                },
                                "null"
                            ]
                        },
                        {
                            "name": "OneTimeUse",
                            "type": [
                                {
                                    "type": "record",
                                    "name": "OneTimeUseType"
                                },
                                "null"
                            ]
                        },
                        {
                            "name": "ProxyRestriction",
                            "type": [
                                {
                                    "type": "record",
                                    "name": "ProxyRestrictionType",
                                    "fields": [
                                        {
                                            "name": "Count",
                                            "type": [
                                                "string",
                                                "null"
                                            ]
                                        },
                                        {
                                            "name": "Audience",
                                            "type": {
                                                "type": "array",
                                                "items": "string"
                                            }
                                        }
                                    ]
                                },
                                "null"
                            ]
                        }
                    ]
                },
                "null"
            ]
        },
        {
            "name": "RequestedAuthnContext",
            "type": [
                {
                    "type": "record",
                    "name": "RequestedAuthnContextType",
                    "fields": [
                        {
                            "name": "Comparison",
                            "type": [
                                "string",
                                "null"
                            ]
                        },
                        {
                            "name": "AuthnContextClassRef",
                            "type": {
                                "type": "array",
                                "items": "string"
                            }
                        },
                        {
                            "name": "AuthnContextDeclRef",
                            "type": {
                                "type": "array",
                                "items": "string"
                            }
                        }
                    ]
                },
                "null"
            ]
        },
        {
            "name": "Scoping",
            "type": [
                {
                    "type": "record",
                    "name": "ScopingType",
                    "fields": [
                        {
                            "name": "ProxyCount",
                            "type": [
                                "string",
                                "null"
                            ]
                        },
                        {
                            "name": "IDPList",
                            "type": [
                                {
                                    "type": "record",
                                    "name": "IDPListType",
                                    "fields": [
                                        {
                                            "name": "IDPEntry",
                                            "type": {
                                                "type": "array",
                                                "items": {
                                                    "type": "record",
                                                    "name": "IDPEntryType",
                                                    "fields": [
                                                        {
                                                            "name": "ProviderID",
                                                            "type": "string"
                                                        },
                                                        {
                                                            "name": "Name",
                                                            "type": [
                                                                "string",
                                                                "null"
                                                            ]
                                                        },
                                                        {
                                                            "name": "Loc",
                                                            "type": [
                                                                "string",
                                                                "null"
                                                            ]
                                                        }
                                                    ]
                                                }
                                            }
                                        },
                                        {
                                            "name": "GetComplete",
                                            "type": [
                                                "string",
                                                "null"
                                            ]
                                        }
                                    ]
                                },
                                "null"
                            ]
                        },
                        {
                            "name": "RequesterID",
                            "type": {
                                "type": "array",
                                "items": "string"
                            }
                        }
                    ]
                },
                "null"
            ]
        }
    ]
}`)

// Generated by codegen. Please do not modify.
var _NameIDType_schema, _NameIDType_schema_err = avro.ParseSchema(`{
    "type": "record",
    "name": "NameIDType",
    "fields": [
        {
            "name": "NameQualifier",
            "type": [
                "string",
                "null"
            ]
        },
        {
            "name": "SPNameQualifier",
            "type": [
                "string",
                "null"
            ]
        },
        {
            "name": "Format",
            "type": [
                "string",
                "null"
            ]
        },
        {
            "name": "SPProvidedID",
            "type": [
                "string",
                "null"
            ]
        }
    ]
}`)

// Generated by codegen. Please do not modify.
var _SignatureType_schema, _SignatureType_schema_err = avro.ParseSchema(`{
    "type": "record",
    "name": "SignatureType",
    "fields": [
        {
            "name": "Id",
            "type": [
                "string",
                "null"
            ]
        },
        {
            "name": "SignedInfo",
            "type": {
                "type": "record",
                "name": "SignedInfoType",
                "fields": [
                    {
                        "name": "Id",
                        "type": [
                            "string",
                            "null"
                        ]
                    },
                    {
                        "name": "CanonicalizationMethod",
                        "type": {
                            "type": "record",
                            "name": "CanonicalizationMethodType",
                            "fields": [
                                {
                                    "name": "Algorithm",
                                    "type": "string"
                                },
                                {
                                    "name": "others",
                                    "type": {
                                        "type": "map",
                                        "values": "string"
                                    }
                                }
                            ]
                        }
                    },
                    {
                        "name": "SignatureMethod",
                        "type": {
                            "type": "record",
                            "name": "SignatureMethodType",
                            "fields": [
                                {
                                    "name": "Algorithm",
                                    "type": "string"
                                },
                                {
                                    "name": "HMACOutputLength",
                                    "type": [
                                        "string",
                                        "null"
                                    ]
                                },
                                {
                                    "name": "others",
                                    "type": {
                                        "type": "map",
                                        "values": "string"
                                    }
                                }
                            ]
                        }
                    },
                    {
                        "name": "Reference",
                        "type": {
                            "type": "array",
                            "items": {
                                "type": "record",
                                "name": "ReferenceType",
                                "fields": [
                                    {
                                        "name": "Id",
                                        "type": [
                                            "string",
                                            "null"
                                        ]
                                    },
                                    {
                                        "name": "URI",
                                        "type": [
                                            "string",
                                            "null"
                                        ]
                                    },
                                    {
                                        "name": "Type",
                                        "type": [
                                            "string",
                                            "null"
                                        ]
                                    },
                                    {
                                        "name": "Transforms",
                                        "type": [
                                            {
                                                "type": "record",
                                                "name": "TransformsType",
                                                "fields": [
                                                    {
                                                        "name": "Transform",
                                                        "type": {
                                                            "type": "array",
                                                            "items": {
                                                                "type": "record",
                                                                "name": "TransformType",
                                                                "fields": [
                                                                    {
                                                                        "name": "Algorithm",
                                                                        "type": "string"
                                                                    },
                                                                    {
                                                                        "name": "others",
                                                                        "type": {
                                                                            "type": "map",
                                                                            "values": "string"
                                                                        }
                                                                    },
                                                                    {
                                                                        "name": "XPath",
                                                                        "type": [
                                                                            "string",
                                                                            "null"
                                                                        ]
                                                                    }
                                                                ]
                                                            }
                                                        }
                                                    }
                                                ]
                                            },
                                            "null"
                                        ]
                                    },
                                    {
                                        "name": "DigestMethod",
                                        "type": {
                                            "type": "record",
                                            "name": "DigestMethodType",
                                            "fields": [
                                                {
                                                    "name": "Algorithm",
                                                    "type": "string"
                                                },
                                                {
                                                    "name": "others",
                                                    "type": {
                                                        "type": "map",
                                                        "values": "string"
                                                    }
                                                }
                                            ]
                                        }
                                    },
                                    {
                                        "name": "DigestValue",
                                        "type": "string"
                                    }
                                ]
                            }
                        }
                    }
                ]
            }
        },
        {
            "name": "SignatureValue",
            "type": {
                "type": "record",
                "name": "SignatureValueType",
                "fields": [
                    {
                        "name": "Id",
                        "type": [
                            "string",
                            "null"
                        ]
                    }
                ]
            }
        },
        {
            "name": "KeyInfo",
            "type": [
                {
                    "type": "record",
                    "name": "KeyInfoType",
                    "fields": [
                        {
                            "name": "Id",
                            "type": [
                                "string",
                                "null"
                            ]
                        },
                        {
                            "name": "KeyName",
                            "type": [
                                "string",
                                "null"
                            ]
                        },
                        {
                            "name": "KeyValue",
                            "type": [
                                {
                                    "type": "record",
                                    "name": "KeyValueType",
                                    "fields": [
                                        {
                                            "name": "DSAKeyValue",
                                            "type": [
                                                {
                                                    "type": "record",
                                                    "name": "DSAKeyValueType",
                                                    "fields": [
                                                        {
                                                            "name": "P",
                                                            "type": "string"
                                                        },
                                                        {
                                                            "name": "Q",
                                                            "type": "string"
                                                        },
                                                        {
                                                            "name": "G",
                                                            "type": [
                                                                "string",
                                                                "null"
                                                            ]
                                                        },
                                                        {
                                                            "name": "Y",
                                                            "type": "string"
                                                        },
                                                        {
                                                            "name": "J",
                                                            "type": [
                                                                "string",
                                                                "null"
                                                            ]
                                                        },
                                                        {
                                                            "name": "Seed",
                                                            "type": "string"
                                                        },
                                                        {
                                                            "name": "PgenCounter",
                                                            "type": "string"
                                                        }
                                                    ]
                                                },
                                                "null"
                                            ]
                                        },
                                        {
                                            "name": "RSAKeyValue",
                                            "type": [
                                                {
                                                    "type": "record",
                                                    "name": "RSAKeyValueType",
                                                    "fields": [
                                                        {
                                                            "name": "Modulus",
                                                            "type": "string"
                                                        },
                                                        {
                                                            "name": "Exponent",
                                                            "type": "string"
                                                        }
                                                    ]
                                                },
                                                "null"
                                            ]
                                        },
                                        {
                                            "name": "others",
                                            "type": {
                                                "type": "map",
                                                "values": "string"
                                            }
                                        }
                                    ]
                                },
                                "null"
                            ]
                        },
                        {
                            "name": "RetrievalMethod",
                            "type": [
                                {
                                    "type": "record",
                                    "name": "RetrievalMethodType",
                                    "fields": [
                                        {
                                            "name": "URI",
                                            "type": [
                                                "string",
                                                "null"
                                            ]
                                        },
                                        {
                                            "name": "Type",
                                            "type": [
                                                "string",
                                                "null"
                                            ]
                                        },
                                        {
                                            "name": "Transforms",
                                            "type": [
                                                "TransformsType",
                                                "null"
                                            ]
                                        }
                                    ]
                                },
                                "null"
                            ]
                        },
                        {
                            "name": "X509Data",
                            "type": [
                                {
                                    "type": "record",
                                    "name": "X509DataType",
                                    "fields": [
                                        {
                                            "name": "X509IssuerSerial",
                                            "type": [
                                                {
                                                    "type": "record",
                                                    "name": "X509IssuerSerialType",
                                                    "fields": [
                                                        {
                                                            "name": "X509IssuerName",
                                                            "type": "string"
                                                        },
                                                        {
                                                            "name": "X509SerialNumber",
                                                            "type": "string"
                                                        }
                                                    ]
                                                },
                                                "null"
                                            ]
                                        },
                                        {
                                            "name": "X509SKI",
                                            "type": [
                                                "string",
                                                "null"
                                            ]
                                        },
                                        {
                                            "name": "X509SubjectName",
                                            "type": [
                                                "string",
                                                "null"
                                            ]
                                        },
                                        {
                                            "name": "X509Certificate",
                                            "type": [
                                                "string",
                                                "null"
                                            ]
                                        },
                                        {
                                            "name": "X509CRL",
                                            "type": [
                                                "string",
                                                "null"
                                            ]
                                        },
                                        {
                                            "name": "others",
                                            "type": {
                                                "type": "map",
                                                "values": "string"
                                            }
                                        }
                                    ]
                                },
                                "null"
                            ]
                        },
                        {
                            "name": "PGPData",
                            "type": [
                                {
                                    "type": "record",
                                    "name": "PGPDataType",
                                    "fields": [
                                        {
                                            "name": "PGPKeyID",
                                            "type": [
                                                "string",
                                                "null"
                                            ]
                                        },
                                        {
                                            "name": "PGPKeyPacket0",
                                            "type": [
                                                "string",
                                                "null"
                                            ]
                                        },
                                        {
                                            "name": "others",
                                            "type": {
                                                "type": "map",
                                                "values": "string"
                                            }
                                        }
                                    ]
                                },
                                "null"
                            ]
                        },
                        {
                            "name": "SPKIData",
                            "type": [
                                {
                                    "type": "record",
                                    "name": "SPKIDataType",
                                    "fields": [
                                        {
                                            "name": "SPKISexp",
                                            "type": "string"
                                        },
                                        {
                                            "name": "others",
                                            "type": {
                                                "type": "map",
                                                "values": "string"
                                            }
                                        }
                                    ]
                                },
                                "null"
                            ]
                        },
                        {
                            "name": "MgmtData",
                            "type": [
                                "string",
                                "null"
                            ]
                        },
                        {
                            "name": "others",
                            "type": {
                                "type": "map",
                                "values": "string"
                            }
                        }
                    ]
                },
                "null"
            ]
        },
        {
            "name": "Object",
            "type": {
                "type": "array",
                "items": {
                    "type": "record",
                    "name": "ObjectType",
                    "fields": [
                        {
                            "name": "Id",
                            "type": [
                                "string",
                                "null"
                            ]
                        },
                        {
                            "name": "MimeType",
                            "type": [
                                "string",
                                "null"
                            ]
                        },
                        {
                            "name": "Encoding",
                            "type": [
                                "string",
                                "null"
                            ]
                        },
                        {
                            "name": "others",
                            "type": {
                                "type": "map",
                                "values": "string"
                            }
                        }
                    ]
                }
            }
        }
    ]
}`)

// Generated by codegen. Please do not modify.
var _SignedInfoType_schema, _SignedInfoType_schema_err = avro.ParseSchema(`{
    "type": "record",
    "name": "SignedInfoType",
    "fields": [
        {
            "name": "Id",
            "type": [
                "string",
                "null"
            ]
        },
        {
            "name": "CanonicalizationMethod",
            "type": {
                "type": "record",
                "name": "CanonicalizationMethodType",
                "fields": [
                    {
                        "name": "Algorithm",
                        "type": "string"
                    },
                    {
                        "name": "others",
                        "type": {
                            "type": "map",
                            "values": "string"
                        }
                    }
                ]
            }
        },
        {
            "name": "SignatureMethod",
            "type": {
                "type": "record",
                "name": "SignatureMethodType",
                "fields": [
                    {
                        "name": "Algorithm",
                        "type": "string"
                    },
                    {
                        "name": "HMACOutputLength",
                        "type": [
                            "string",
                            "null"
                        ]
                    },
                    {
                        "name": "others",
                        "type": {
                            "type": "map",
                            "values": "string"
                        }
                    }
                ]
            }
        },
        {
            "name": "Reference",
            "type": {
                "type": "array",
                "items": {
                    "type": "record",
                    "name": "ReferenceType",
                    "fields": [
                        {
                            "name": "Id",
                            "type": [
                                "string",
                                "null"
                            ]
                        },
                        {
                            "name": "URI",
                            "type": [
                                "string",
                                "null"
                            ]
                        },
                        {
                            "name": "Type",
                            "type": [
                                "string",
                                "null"
                            ]
                        },
                        {
                            "name": "Transforms",
                            "type": [
                                {
                                    "type": "record",
                                    "name": "TransformsType",
                                    "fields": [
                                        {
                                            "name": "Transform",
                                            "type": {
                                                "type": "array",
                                                "items": {
                                                    "type": "record",
                                                    "name": "TransformType",
                                                    "fields": [
                                                        {
                                                            "name": "Algorithm",
                                                            "type": "string"
                                                        },
                                                        {
                                                            "name": "others",
                                                            "type": {
                                                                "type": "map",
                                                                "values": "string"
                                                            }
                                                        },
                                                        {
                                                            "name": "XPath",
                                                            "type": [
                                                                "string",
                                                                "null"
                                                            ]
                                                        }
                                                    ]
                                                }
                                            }
                                        }
                                    ]
                                },
                                "null"
                            ]
                        },
                        {
                            "name": "DigestMethod",
                            "type": {
                                "type": "record",
                                "name": "DigestMethodType",
                                "fields": [
                                    {
                                        "name": "Algorithm",
                                        "type": "string"
                                    },
                                    {
                                        "name": "others",
                                        "type": {
                                            "type": "map",
                                            "values": "string"
                                        }
                                    }
                                ]
                            }
                        },
                        {
                            "name": "DigestValue",
                            "type": "string"
                        }
                    ]
                }
            }
        }
    ]
}`)

// Generated by codegen. Please do not modify.
var _CanonicalizationMethodType_schema, _CanonicalizationMethodType_schema_err = avro.ParseSchema(`{
    "type": "record",
    "name": "CanonicalizationMethodType",
    "fields": [
        {
            "name": "Algorithm",
            "type": "string"
        },
        {
            "name": "others",
            "type": {
                "type": "map",
                "values": "string"
            }
        }
    ]
}`)

// Generated by codegen. Please do not modify.
var _SignatureMethodType_schema, _SignatureMethodType_schema_err = avro.ParseSchema(`{
    "type": "record",
    "name": "SignatureMethodType",
    "fields": [
        {
            "name": "Algorithm",
            "type": "string"
        },
        {
            "name": "HMACOutputLength",
            "type": [
                "string",
                "null"
            ]
        },
        {
            "name": "others",
            "type": {
                "type": "map",
                "values": "string"
            }
        }
    ]
}`)

// Generated by codegen. Please do not modify.
var _ReferenceType_schema, _ReferenceType_schema_err = avro.ParseSchema(`{
    "type": "record",
    "name": "ReferenceType",
    "fields": [
        {
            "name": "Id",
            "type": [
                "string",
                "null"
            ]
        },
        {
            "name": "URI",
            "type": [
                "string",
                "null"
            ]
        },
        {
            "name": "Type",
            "type": [
                "string",
                "null"
            ]
        },
        {
            "name": "Transforms",
            "type": [
                {
                    "type": "record",
                    "name": "TransformsType",
                    "fields": [
                        {
                            "name": "Transform",
                            "type": {
                                "type": "array",
                                "items": {
                                    "type": "record",
                                    "name": "TransformType",
                                    "fields": [
                                        {
                                            "name": "Algorithm",
                                            "type": "string"
                                        },
                                        {
                                            "name": "others",
                                            "type": {
                                                "type": "map",
                                                "values": "string"
                                            }
                                        },
                                        {
                                            "name": "XPath",
                                            "type": [
                                                "string",
                                                "null"
                                            ]
                                        }
                                    ]
                                }
                            }
                        }
                    ]
                },
                "null"
            ]
        },
        {
            "name": "DigestMethod",
            "type": {
                "type": "record",
                "name": "DigestMethodType",
                "fields": [
                    {
                        "name": "Algorithm",
                        "type": "string"
                    },
                    {
                        "name": "others",
                        "type": {
                            "type": "map",
                            "values": "string"
                        }
                    }
                ]
            }
        },
        {
            "name": "DigestValue",
            "type": "string"
        }
    ]
}`)

// Generated by codegen. Please do not modify.
var _TransformsType_schema, _TransformsType_schema_err = avro.ParseSchema(`{
    "type": "record",
    "name": "TransformsType",
    "fields": [
        {
            "name": "Transform",
            "type": {
                "type": "array",
                "items": {
                    "type": "record",
                    "name": "TransformType",
                    "fields": [
                        {
                            "name": "Algorithm",
                            "type": "string"
                        },
                        {
                            "name": "others",
                            "type": {
                                "type": "map",
                                "values": "string"
                            }
                        },
                        {
                            "name": "XPath",
                            "type": [
                                "string",
                                "null"
                            ]
                        }
                    ]
                }
            }
        }
    ]
}`)

// Generated by codegen. Please do not modify.
var _TransformType_schema, _TransformType_schema_err = avro.ParseSchema(`{
    "type": "record",
    "name": "TransformType",
    "fields": [
        {
            "name": "Algorithm",
            "type": "string"
        },
        {
            "name": "others",
            "type": {
                "type": "map",
                "values": "string"
            }
        },
        {
            "name": "XPath",
            "type": [
                "string",
                "null"
            ]
        }
    ]
}`)

// Generated by codegen. Please do not modify.
var _DigestMethodType_schema, _DigestMethodType_schema_err = avro.ParseSchema(`{
    "type": "record",
    "name": "DigestMethodType",
    "fields": [
        {
            "name": "Algorithm",
            "type": "string"
        },
        {
            "name": "others",
            "type": {
                "type": "map",
                "values": "string"
            }
        }
    ]
}`)

// Generated by codegen. Please do not modify.
var _SignatureValueType_schema, _SignatureValueType_schema_err = avro.ParseSchema(`{
    "type": "record",
    "name": "SignatureValueType",
    "fields": [
        {
            "name": "Id",
            "type": [
                "string",
                "null"
            ]
        }
    ]
}`)

// Generated by codegen. Please do not modify.
var _KeyInfoType_schema, _KeyInfoType_schema_err = avro.ParseSchema(`{
    "type": "record",
    "name": "KeyInfoType",
    "fields": [
        {
            "name": "Id",
            "type": [
                "string",
                "null"
            ]
        },
        {
            "name": "KeyName",
            "type": [
                "string",
                "null"
            ]
        },
        {
            "name": "KeyValue",
            "type": [
                {
                    "type": "record",
                    "name": "KeyValueType",
                    "fields": [
                        {
                            "name": "DSAKeyValue",
                            "type": [
                                {
                                    "type": "record",
                                    "name": "DSAKeyValueType",
                                    "fields": [
                                        {
                                            "name": "P",
                                            "type": "string"
                                        },
                                        {
                                            "name": "Q",
                                            "type": "string"
                                        },
                                        {
                                            "name": "G",
                                            "type": [
                                                "string",
                                                "null"
                                            ]
                                        },
                                        {
                                            "name": "Y",
                                            "type": "string"
                                        },
                                        {
                                            "name": "J",
                                            "type": [
                                                "string",
                                                "null"
                                            ]
                                        },
                                        {
                                            "name": "Seed",
                                            "type": "string"
                                        },
                                        {
                                            "name": "PgenCounter",
                                            "type": "string"
                                        }
                                    ]
                                },
                                "null"
                            ]
                        },
                        {
                            "name": "RSAKeyValue",
                            "type": [
                                {
                                    "type": "record",
                                    "name": "RSAKeyValueType",
                                    "fields": [
                                        {
                                            "name": "Modulus",
                                            "type": "string"
                                        },
                                        {
                                            "name": "Exponent",
                                            "type": "string"
                                        }
                                    ]
                                },
                                "null"
                            ]
                        },
                        {
                            "name": "others",
                            "type": {
                                "type": "map",
                                "values": "string"
                            }
                        }
                    ]
                },
                "null"
            ]
        },
        {
            "name": "RetrievalMethod",
            "type": [
                {
                    "type": "record",
                    "name": "RetrievalMethodType",
                    "fields": [
                        {
                            "name": "URI",
                            "type": [
                                "string",
                                "null"
                            ]
                        },
                        {
                            "name": "Type",
                            "type": [
                                "string",
                                "null"
                            ]
                        },
                        {
                            "name": "Transforms",
                            "type": [
                                "TransformsType",
                                "null"
                            ]
                        }
                    ]
                },
                "null"
            ]
        },
        {
            "name": "X509Data",
            "type": [
                {
                    "type": "record",
                    "name": "X509DataType",
                    "fields": [
                        {
                            "name": "X509IssuerSerial",
                            "type": [
                                {
                                    "type": "record",
                                    "name": "X509IssuerSerialType",
                                    "fields": [
                                        {
                                            "name": "X509IssuerName",
                                            "type": "string"
                                        },
                                        {
                                            "name": "X509SerialNumber",
                                            "type": "string"
                                        }
                                    ]
                                },
                                "null"
                            ]
                        },
                        {
                            "name": "X509SKI",
                            "type": [
                                "string",
                                "null"
                            ]
                        },
                        {
                            "name": "X509SubjectName",
                            "type": [
                                "string",
                                "null"
                            ]
                        },
                        {
                            "name": "X509Certificate",
                            "type": [
                                "string",
                                "null"
                            ]
                        },
                        {
                            "name": "X509CRL",
                            "type": [
                                "string",
                                "null"
                            ]
                        },
                        {
                            "name": "others",
                            "type": {
                                "type": "map",
                                "values": "string"
                            }
                        }
                    ]
                },
                "null"
            ]
        },
        {
            "name": "PGPData",
            "type": [
                {
                    "type": "record",
                    "name": "PGPDataType",
                    "fields": [
                        {
                            "name": "PGPKeyID",
                            "type": [
                                "string",
                                "null"
                            ]
                        },
                        {
                            "name": "PGPKeyPacket0",
                            "type": [
                                "string",
                                "null"
                            ]
                        },
                        {
                            "name": "others",
                            "type": {
                                "type": "map",
                                "values": "string"
                            }
                        }
                    ]
                },
                "null"
            ]
        },
        {
            "name": "SPKIData",
            "type": [
                {
                    "type": "record",
                    "name": "SPKIDataType",
                    "fields": [
                        {
                            "name": "SPKISexp",
                            "type": "string"
                        },
                        {
                            "name": "others",
                            "type": {
                                "type": "map",
                                "values": "string"
                            }
                        }
                    ]
                },
                "null"
            ]
        },
        {
            "name": "MgmtData",
            "type": [
                "string",
                "null"
            ]
        },
        {
            "name": "others",
            "type": {
                "type": "map",
                "values": "string"
            }
        }
    ]
}`)

// Generated by codegen. Please do not modify.
var _KeyValueType_schema, _KeyValueType_schema_err = avro.ParseSchema(`{
    "type": "record",
    "name": "KeyValueType",
    "fields": [
        {
            "name": "DSAKeyValue",
            "type": [
                {
                    "type": "record",
                    "name": "DSAKeyValueType",
                    "fields": [
                        {
                            "name": "P",
                            "type": "string"
                        },
                        {
                            "name": "Q",
                            "type": "string"
                        },
                        {
                            "name": "G",
                            "type": [
                                "string",
                                "null"
                            ]
                        },
                        {
                            "name": "Y",
                            "type": "string"
                        },
                        {
                            "name": "J",
                            "type": [
                                "string",
                                "null"
                            ]
                        },
                        {
                            "name": "Seed",
                            "type": "string"
                        },
                        {
                            "name": "PgenCounter",
                            "type": "string"
                        }
                    ]
                },
                "null"
            ]
        },
        {
            "name": "RSAKeyValue",
            "type": [
                {
                    "type": "record",
                    "name": "RSAKeyValueType",
                    "fields": [
                        {
                            "name": "Modulus",
                            "type": "string"
                        },
                        {
                            "name": "Exponent",
                            "type": "string"
                        }
                    ]
                },
                "null"
            ]
        },
        {
            "name": "others",
            "type": {
                "type": "map",
                "values": "string"
            }
        }
    ]
}`)

// Generated by codegen. Please do not modify.
var _DSAKeyValueType_schema, _DSAKeyValueType_schema_err = avro.ParseSchema(`{
    "type": "record",
    "name": "DSAKeyValueType",
    "fields": [
        {
            "name": "P",
            "type": "string"
        },
        {
            "name": "Q",
            "type": "string"
        },
        {
            "name": "G",
            "type": [
                "string",
                "null"
            ]
        },
        {
            "name": "Y",
            "type": "string"
        },
        {
            "name": "J",
            "type": [
                "string",
                "null"
            ]
        },
        {
            "name": "Seed",
            "type": "string"
        },
        {
            "name": "PgenCounter",
            "type": "string"
        }
    ]
}`)

// Generated by codegen. Please do not modify.
var _RSAKeyValueType_schema, _RSAKeyValueType_schema_err = avro.ParseSchema(`{
    "type": "record",
    "name": "RSAKeyValueType",
    "fields": [
        {
            "name": "Modulus",
            "type": "string"
        },
        {
            "name": "Exponent",
            "type": "string"
        }
    ]
}`)

// Generated by codegen. Please do not modify.
var _RetrievalMethodType_schema, _RetrievalMethodType_schema_err = avro.ParseSchema(`{
    "type": "record",
    "name": "RetrievalMethodType",
    "fields": [
        {
            "name": "URI",
            "type": [
                "string",
                "null"
            ]
        },
        {
            "name": "Type",
            "type": [
                "string",
                "null"
            ]
        },
        {
            "name": "Transforms",
            "type": [
                "TransformsType",
                "null"
            ]
        }
    ]
}`)

// Generated by codegen. Please do not modify.
var _X509DataType_schema, _X509DataType_schema_err = avro.ParseSchema(`{
    "type": "record",
    "name": "X509DataType",
    "fields": [
        {
            "name": "X509IssuerSerial",
            "type": [
                {
                    "type": "record",
                    "name": "X509IssuerSerialType",
                    "fields": [
                        {
                            "name": "X509IssuerName",
                            "type": "string"
                        },
                        {
                            "name": "X509SerialNumber",
                            "type": "string"
                        }
                    ]
                },
                "null"
            ]
        },
        {
            "name": "X509SKI",
            "type": [
                "string",
                "null"
            ]
        },
        {
            "name": "X509SubjectName",
            "type": [
                "string",
                "null"
            ]
        },
        {
            "name": "X509Certificate",
            "type": [
                "string",
                "null"
            ]
        },
        {
            "name": "X509CRL",
            "type": [
                "string",
                "null"
            ]
        },
        {
            "name": "others",
            "type": {
                "type": "map",
                "values": "string"
            }
        }
    ]
}`)

// Generated by codegen. Please do not modify.
var _X509IssuerSerialType_schema, _X509IssuerSerialType_schema_err = avro.ParseSchema(`{
    "type": "record",
    "name": "X509IssuerSerialType",
    "fields": [
        {
            "name": "X509IssuerName",
            "type": "string"
        },
        {
            "name": "X509SerialNumber",
            "type": "string"
        }
    ]
}`)

// Generated by codegen. Please do not modify.
var _PGPDataType_schema, _PGPDataType_schema_err = avro.ParseSchema(`{
    "type": "record",
    "name": "PGPDataType",
    "fields": [
        {
            "name": "PGPKeyID",
            "type": [
                "string",
                "null"
            ]
        },
        {
            "name": "PGPKeyPacket0",
            "type": [
                "string",
                "null"
            ]
        },
        {
            "name": "others",
            "type": {
                "type": "map",
                "values": "string"
            }
        }
    ]
}`)

// Generated by codegen. Please do not modify.
var _SPKIDataType_schema, _SPKIDataType_schema_err = avro.ParseSchema(`{
    "type": "record",
    "name": "SPKIDataType",
    "fields": [
        {
            "name": "SPKISexp",
            "type": "string"
        },
        {
            "name": "others",
            "type": {
                "type": "map",
                "values": "string"
            }
        }
    ]
}`)

// Generated by codegen. Please do not modify.
var _ObjectType_schema, _ObjectType_schema_err = avro.ParseSchema(`{
    "type": "record",
    "name": "ObjectType",
    "fields": [
        {
            "name": "Id",
            "type": [
                "string",
                "null"
            ]
        },
        {
            "name": "MimeType",
            "type": [
                "string",
                "null"
            ]
        },
        {
            "name": "Encoding",
            "type": [
                "string",
                "null"
            ]
        },
        {
            "name": "others",
            "type": {
                "type": "map",
                "values": "string"
            }
        }
    ]
}`)

// Generated by codegen. Please do not modify.
var _ExtensionsType_schema, _ExtensionsType_schema_err = avro.ParseSchema(`{
    "type": "record",
    "name": "ExtensionsType",
    "fields": [
        {
            "name": "others",
            "type": {
                "type": "map",
                "values": "string"
            }
        }
    ]
}`)

// Generated by codegen. Please do not modify.
var _SubjectType_schema, _SubjectType_schema_err = avro.ParseSchema(`{
    "type": "record",
    "name": "SubjectType",
    "fields": [
        {
            "name": "BaseID",
            "type": [
                {
                    "type": "record",
                    "name": "BaseIDAbstractType",
                    "fields": [
                        {
                            "name": "NameQualifier",
                            "type": [
                                "string",
                                "null"
                            ]
                        },
                        {
                            "name": "SPNameQualifier",
                            "type": [
                                "string",
                                "null"
                            ]
                        }
                    ]
                },
                "null"
            ]
        },
        {
            "name": "NameID",
            "type": [
                "NameIDType",
                "null"
            ]
        },
        {
            "name": "EncryptedID",
            "type": [
                {
                    "type": "record",
                    "name": "EncryptedElementType",
                    "fields": [
                        {
                            "name": "EncryptedData",
                            "type": {
                                "type": "record",
                                "name": "EncryptedDataType",
                                "fields": [
                                    {
                                        "name": "Id",
                                        "type": [
                                            "string",
                                            "null"
                                        ]
                                    },
                                    {
                                        "name": "Type",
                                        "type": [
                                            "string",
                                            "null"
                                        ]
                                    },
                                    {
                                        "name": "MimeType",
                                        "type": [
                                            "string",
                                            "null"
                                        ]
                                    },
                                    {
                                        "name": "Encoding",
                                        "type": [
                                            "string",
                                            "null"
                                        ]
                                    },
                                    {
                                        "name": "EncryptionMethod",
                                        "type": [
                                            {
                                                "type": "record",
                                                "name": "EncryptionMethodType",
                                                "fields": [
                                                    {
                                                        "name": "Algorithm",
                                                        "type": "string"
                                                    },
                                                    {
                                                        "name": "KeySize",
                                                        "type": [
                                                            "string",
                                                            "null"
                                                        ]
                                                    },
                                                    {
                                                        "name": "OAEPparams",
                                                        "type": [
                                                            "string",
                                                            "null"
                                                        ]
                                                    },
                                                    {
                                                        "name": "others",
                                                        "type": {
                                                            "type": "map",
                                                            "values": "string"
                                                        }
                                                    }
                                                ]
                                            },
                                            "null"
                                        ]
                                    },
                                    {
                                        "name": "KeyInfo",
                                        "type": [
                                            "KeyInfoType",
                                            "null"
                                        ]
                                    },
                                    {
                                        "name": "CipherData",
                                        "type": {
                                            "type": "record",
                                            "name": "CipherDataType",
                                            "fields": [
                                                {
                                                    "name": "CipherValue",
                                                    "type": [
                                                        "string",
                                                        "null"
                                                    ]
                                                },
                                                {
                                                    "name": "CipherReference",
                                                    "type": [
                                                        {
                                                            "type": "record",
                                                            "name": "CipherReferenceType",
                                                            "fields": [
                                                                {
                                                                    "name": "URI",
                                                                    "type": "string"
                                                                },
                                                                {
                                                                    "name": "Transforms",
                                                                    "type": [
                                                                        "TransformsType",
                                                                        "null"
                                                                    ]
                                                                }
                                                            ]
                                                        },
                                                        "null"
                                                    ]
                                                }
                                            ]
                                        }
                                    },
                                    {
                                        "name": "EncryptionProperties",
                                        "type": [
                                            {
                                                "type": "record",
                                                "name": "EncryptionPropertiesType",
                                                "fields": [
                                                    {
                                                        "name": "Id",
                                                        "type": [
                                                            "string",
                                                            "null"
                                                        ]
                                                    },
                                                    {
                                                        "name": "EncryptionProperty",
                                                        "type": {
                                                            "type": "array",
                                                            "items": {
                                                                "type": "record",
                                                                "name": "EncryptionPropertyType",
                                                                "fields": [
                                                                    {
                                                                        "name": "Target",
                                                                        "type": [
                                                                            "string",
                                                                            "null"
                                                                        ]
                                                                    },
                                                                    {
                                                                        "name": "Id",
                                                                        "type": [
                                                                            "string",
                                                                            "null"
                                                                        ]
                                                                    },
                                                                    {
                                                                        "name": "others",
                                                                        "type": {
                                                                            "type": "map",
                                                                            "values": "string"
                                                                        }
                                                                    }
                                                                ]
                                                            }
                                                        }
                                                    }
                                                ]
                                            },
                                            "null"
                                        ]
                                    }
                                ]
                            }
                        },
                        {
                            "name": "EncryptedKey",
                            "type": {
                                "type": "array",
                                "items": {
                                    "type": "record",
                                    "name": "EncryptedKeyType",
                                    "fields": [
                                        {
                                            "name": "Recipient",
                                            "type": [
                                                "string",
                                                "null"
                                            ]
                                        },
                                        {
                                            "name": "Id",
                                            "type": [
                                                "string",
                                                "null"
                                            ]
                                        },
                                        {
                                            "name": "Type",
                                            "type": [
                                                "string",
                                                "null"
                                            ]
                                        },
                                        {
                                            "name": "MimeType",
                                            "type": [
                                                "string",
                                                "null"
                                            ]
                                        },
                                        {
                                            "name": "Encoding",
                                            "type": [
                                                "string",
                                                "null"
                                            ]
                                        },
                                        {
                                            "name": "EncryptionMethod",
                                            "type": [
                                                "EncryptionMethodType",
                                                "null"
                                            ]
                                        },
                                        {
                                            "name": "KeyInfo",
                                            "type": [
                                                "KeyInfoType",
                                                "null"
                                            ]
                                        },
                                        {
                                            "name": "CipherData",
                                            "type": "CipherDataType"
                                        },
                                        {
                                            "name": "EncryptionProperties",
                                            "type": [
                                                "EncryptionPropertiesType",
                                                "null"
                                            ]
                                        },
                                        {
                                            "name": "ReferenceList",
                                            "type": [
                                                {
                                                    "type": "record",
                                                    "name": "type0",
                                                    "fields": [
                                                        {
                                                            "name": "DataReference",
                                                            "type": [
                                                                "ReferenceType",
                                                                "null"
                                                            ]
                                                        },
                                                        {
                                                            "name": "KeyReference",
                                                            "type": [
                                                                "ReferenceType",
                                                                "null"
                                                            ]
                                                        }
                                                    ]
                                                },
                                                "null"
                                            ]
                                        },
                                        {
                                            "name": "CarriedKeyName",
                                            "type": [
                                                "string",
                                                "null"
                                            ]
                                        }
                                    ]
                                }
                            }
                        }
                    ]
                },
                "null"
            ]
        },
        {
            "name": "SubjectConfirmation0",
            "type": {
                "type": "array",
                "items": {
                    "type": "record",
                    "name": "SubjectConfirmationType",
                    "fields": [
                        {
                            "name": "Method",
                            "type": "string"
                        },
                        {
                            "name": "BaseID",
                            "type": [
                                "BaseIDAbstractType",
                                "null"
                            ]
                        },
                        {
                            "name": "NameID",
                            "type": [
                                "NameIDType",
                                "null"
                            ]
                        },
                        {
                            "name": "EncryptedID",
                            "type": [
                                "EncryptedElementType",
                                "null"
                            ]
                        },
                        {
                            "name": "SubjectConfirmationData",
                            "type": [
                                {
                                    "type": "record",
                                    "name": "SubjectConfirmationDataType",
                                    "fields": [
                                        {
                                            "name": "NotBefore",
                                            "type": [
                                                "string",
                                                "null"
                                            ]
                                        },
                                        {
                                            "name": "NotOnOrAfter",
                                            "type": [
                                                "string",
                                                "null"
                                            ]
                                        },
                                        {
                                            "name": "Recipient",
                                            "type": [
                                                "string",
                                                "null"
                                            ]
                                        },
                                        {
                                            "name": "InResponseTo",
                                            "type": [
                                                "string",
                                                "null"
                                            ]
                                        },
                                        {
                                            "name": "Address",
                                            "type": [
                                                "string",
                                                "null"
                                            ]
                                        },
                                        {
                                            "name": "others",
                                            "type": {
                                                "type": "map",
                                                "values": "string"
                                            }
                                        }
                                    ]
                                },
                                "null"
                            ]
                        }
                    ]
                }
            }
        }
    ]
}`)

// Generated by codegen. Please do not modify.
var _BaseIDAbstractType_schema, _BaseIDAbstractType_schema_err = avro.ParseSchema(`{
    "type": "record",
    "name": "BaseIDAbstractType",
    "fields": [
        {
            "name": "NameQualifier",
            "type": [
                "string",
                "null"
            ]
        },
        {
            "name": "SPNameQualifier",
            "type": [
                "string",
                "null"
            ]
        }
    ]
}`)

// Generated by codegen. Please do not modify.
var _EncryptedElementType_schema, _EncryptedElementType_schema_err = avro.ParseSchema(`{
    "type": "record",
    "name": "EncryptedElementType",
    "fields": [
        {
            "name": "EncryptedData",
            "type": {
                "type": "record",
                "name": "EncryptedDataType",
                "fields": [
                    {
                        "name": "Id",
                        "type": [
                            "string",
                            "null"
                        ]
                    },
                    {
                        "name": "Type",
                        "type": [
                            "string",
                            "null"
                        ]
                    },
                    {
                        "name": "MimeType",
                        "type": [
                            "string",
                            "null"
                        ]
                    },
                    {
                        "name": "Encoding",
                        "type": [
                            "string",
                            "null"
                        ]
                    },
                    {
                        "name": "EncryptionMethod",
                        "type": [
                            {
                                "type": "record",
                                "name": "EncryptionMethodType",
                                "fields": [
                                    {
                                        "name": "Algorithm",
                                        "type": "string"
                                    },
                                    {
                                        "name": "KeySize",
                                        "type": [
                                            "string",
                                            "null"
                                        ]
                                    },
                                    {
                                        "name": "OAEPparams",
                                        "type": [
                                            "string",
                                            "null"
                                        ]
                                    },
                                    {
                                        "name": "others",
                                        "type": {
                                            "type": "map",
                                            "values": "string"
                                        }
                                    }
                                ]
                            },
                            "null"
                        ]
                    },
                    {
                        "name": "KeyInfo",
                        "type": [
                            "KeyInfoType",
                            "null"
                        ]
                    },
                    {
                        "name": "CipherData",
                        "type": {
                            "type": "record",
                            "name": "CipherDataType",
                            "fields": [
                                {
                                    "name": "CipherValue",
                                    "type": [
                                        "string",
                                        "null"
                                    ]
                                },
                                {
                                    "name": "CipherReference",
                                    "type": [
                                        {
                                            "type": "record",
                                            "name": "CipherReferenceType",
                                            "fields": [
                                                {
                                                    "name": "URI",
                                                    "type": "string"
                                                },
                                                {
                                                    "name": "Transforms",
                                                    "type": [
                                                        "TransformsType",
                                                        "null"
                                                    ]
                                                }
                                            ]
                                        },
                                        "null"
                                    ]
                                }
                            ]
                        }
                    },
                    {
                        "name": "EncryptionProperties",
                        "type": [
                            {
                                "type": "record",
                                "name": "EncryptionPropertiesType",
                                "fields": [
                                    {
                                        "name": "Id",
                                        "type": [
                                            "string",
                                            "null"
                                        ]
                                    },
                                    {
                                        "name": "EncryptionProperty",
                                        "type": {
                                            "type": "array",
                                            "items": {
                                                "type": "record",
                                                "name": "EncryptionPropertyType",
                                                "fields": [
                                                    {
                                                        "name": "Target",
                                                        "type": [
                                                            "string",
                                                            "null"
                                                        ]
                                                    },
                                                    {
                                                        "name": "Id",
                                                        "type": [
                                                            "string",
                                                            "null"
                                                        ]
                                                    },
                                                    {
                                                        "name": "others",
                                                        "type": {
                                                            "type": "map",
                                                            "values": "string"
                                                        }
                                                    }
                                                ]
                                            }
                                        }
                                    }
                                ]
                            },
                            "null"
                        ]
                    }
                ]
            }
        },
        {
            "name": "EncryptedKey",
            "type": {
                "type": "array",
                "items": {
                    "type": "record",
                    "name": "EncryptedKeyType",
                    "fields": [
                        {
                            "name": "Recipient",
                            "type": [
                                "string",
                                "null"
                            ]
                        },
                        {
                            "name": "Id",
                            "type": [
                                "string",
                                "null"
                            ]
                        },
                        {
                            "name": "Type",
                            "type": [
                                "string",
                                "null"
                            ]
                        },
                        {
                            "name": "MimeType",
                            "type": [
                                "string",
                                "null"
                            ]
                        },
                        {
                            "name": "Encoding",
                            "type": [
                                "string",
                                "null"
                            ]
                        },
                        {
                            "name": "EncryptionMethod",
                            "type": [
                                "EncryptionMethodType",
                                "null"
                            ]
                        },
                        {
                            "name": "KeyInfo",
                            "type": [
                                "KeyInfoType",
                                "null"
                            ]
                        },
                        {
                            "name": "CipherData",
                            "type": "CipherDataType"
                        },
                        {
                            "name": "EncryptionProperties",
                            "type": [
                                "EncryptionPropertiesType",
                                "null"
                            ]
                        },
                        {
                            "name": "ReferenceList",
                            "type": [
                                {
                                    "type": "record",
                                    "name": "type0",
                                    "fields": [
                                        {
                                            "name": "DataReference",
                                            "type": [
                                                "ReferenceType",
                                                "null"
                                            ]
                                        },
                                        {
                                            "name": "KeyReference",
                                            "type": [
                                                "ReferenceType",
                                                "null"
                                            ]
                                        }
                                    ]
                                },
                                "null"
                            ]
                        },
                        {
                            "name": "CarriedKeyName",
                            "type": [
                                "string",
                                "null"
                            ]
                        }
                    ]
                }
            }
        }
    ]
}`)

// Generated by codegen. Please do not modify.
var _EncryptedDataType_schema, _EncryptedDataType_schema_err = avro.ParseSchema(`{
    "type": "record",
    "name": "EncryptedDataType",
    "fields": [
        {
            "name": "Id",
            "type": [
                "string",
                "null"
            ]
        },
        {
            "name": "Type",
            "type": [
                "string",
                "null"
            ]
        },
        {
            "name": "MimeType",
            "type": [
                "string",
                "null"
            ]
        },
        {
            "name": "Encoding",
            "type": [
                "string",
                "null"
            ]
        },
        {
            "name": "EncryptionMethod",
            "type": [
                {
                    "type": "record",
                    "name": "EncryptionMethodType",
                    "fields": [
                        {
                            "name": "Algorithm",
                            "type": "string"
                        },
                        {
                            "name": "KeySize",
                            "type": [
                                "string",
                                "null"
                            ]
                        },
                        {
                            "name": "OAEPparams",
                            "type": [
                                "string",
                                "null"
                            ]
                        },
                        {
                            "name": "others",
                            "type": {
                                "type": "map",
                                "values": "string"
                            }
                        }
                    ]
                },
                "null"
            ]
        },
        {
            "name": "KeyInfo",
            "type": [
                "KeyInfoType",
                "null"
            ]
        },
        {
            "name": "CipherData",
            "type": {
                "type": "record",
                "name": "CipherDataType",
                "fields": [
                    {
                        "name": "CipherValue",
                        "type": [
                            "string",
                            "null"
                        ]
                    },
                    {
                        "name": "CipherReference",
                        "type": [
                            {
                                "type": "record",
                                "name": "CipherReferenceType",
                                "fields": [
                                    {
                                        "name": "URI",
                                        "type": "string"
                                    },
                                    {
                                        "name": "Transforms",
                                        "type": [
                                            "TransformsType",
                                            "null"
                                        ]
                                    }
                                ]
                            },
                            "null"
                        ]
                    }
                ]
            }
        },
        {
            "name": "EncryptionProperties",
            "type": [
                {
                    "type": "record",
                    "name": "EncryptionPropertiesType",
                    "fields": [
                        {
                            "name": "Id",
                            "type": [
                                "string",
                                "null"
                            ]
                        },
                        {
                            "name": "EncryptionProperty",
                            "type": {
                                "type": "array",
                                "items": {
                                    "type": "record",
                                    "name": "EncryptionPropertyType",
                                    "fields": [
                                        {
                                            "name": "Target",
                                            "type": [
                                                "string",
                                                "null"
                                            ]
                                        },
                                        {
                                            "name": "Id",
                                            "type": [
                                                "string",
                                                "null"
                                            ]
                                        },
                                        {
                                            "name": "others",
                                            "type": {
                                                "type": "map",
                                                "values": "string"
                                            }
                                        }
                                    ]
                                }
                            }
                        }
                    ]
                },
                "null"
            ]
        }
    ]
}`)

// Generated by codegen. Please do not modify.
var _EncryptionMethodType_schema, _EncryptionMethodType_schema_err = avro.ParseSchema(`{
    "type": "record",
    "name": "EncryptionMethodType",
    "fields": [
        {
            "name": "Algorithm",
            "type": "string"
        },
        {
            "name": "KeySize",
            "type": [
                "string",
                "null"
            ]
        },
        {
            "name": "OAEPparams",
            "type": [
                "string",
                "null"
            ]
        },
        {
            "name": "others",
            "type": {
                "type": "map",
                "values": "string"
            }
        }
    ]
}`)

// Generated by codegen. Please do not modify.
var _CipherDataType_schema, _CipherDataType_schema_err = avro.ParseSchema(`{
    "type": "record",
    "name": "CipherDataType",
    "fields": [
        {
            "name": "CipherValue",
            "type": [
                "string",
                "null"
            ]
        },
        {
            "name": "CipherReference",
            "type": [
                {
                    "type": "record",
                    "name": "CipherReferenceType",
                    "fields": [
                        {
                            "name": "URI",
                            "type": "string"
                        },
                        {
                            "name": "Transforms",
                            "type": [
                                "TransformsType",
                                "null"
                            ]
                        }
                    ]
                },
                "null"
            ]
        }
    ]
}`)

// Generated by codegen. Please do not modify.
var _CipherReferenceType_schema, _CipherReferenceType_schema_err = avro.ParseSchema(`{
    "type": "record",
    "name": "CipherReferenceType",
    "fields": [
        {
            "name": "URI",
            "type": "string"
        },
        {
            "name": "Transforms",
            "type": [
                "TransformsType",
                "null"
            ]
        }
    ]
}`)

// Generated by codegen. Please do not modify.
var _EncryptionPropertiesType_schema, _EncryptionPropertiesType_schema_err = avro.ParseSchema(`{
    "type": "record",
    "name": "EncryptionPropertiesType",
    "fields": [
        {
            "name": "Id",
            "type": [
                "string",
                "null"
            ]
        },
        {
            "name": "EncryptionProperty",
            "type": {
                "type": "array",
                "items": {
                    "type": "record",
                    "name": "EncryptionPropertyType",
                    "fields": [
                        {
                            "name": "Target",
                            "type": [
                                "string",
                                "null"
                            ]
                        },
                        {
                            "name": "Id",
                            "type": [
                                "string",
                                "null"
                            ]
                        },
                        {
                            "name": "others",
                            "type": {
                                "type": "map",
                                "values": "string"
                            }
                        }
                    ]
                }
            }
        }
    ]
}`)

// Generated by codegen. Please do not modify.
var _EncryptionPropertyType_schema, _EncryptionPropertyType_schema_err = avro.ParseSchema(`{
    "type": "record",
    "name": "EncryptionPropertyType",
    "fields": [
        {
            "name": "Target",
            "type": [
                "string",
                "null"
            ]
        },
        {
            "name": "Id",
            "type": [
                "string",
                "null"
            ]
        },
        {
            "name": "others",
            "type": {
                "type": "map",
                "values": "string"
            }
        }
    ]
}`)

// Generated by codegen. Please do not modify.
var _EncryptedKeyType_schema, _EncryptedKeyType_schema_err = avro.ParseSchema(`{
    "type": "record",
    "name": "EncryptedKeyType",
    "fields": [
        {
            "name": "Recipient",
            "type": [
                "string",
                "null"
            ]
        },
        {
            "name": "Id",
            "type": [
                "string",
                "null"
            ]
        },
        {
            "name": "Type",
            "type": [
                "string",
                "null"
            ]
        },
        {
            "name": "MimeType",
            "type": [
                "string",
                "null"
            ]
        },
        {
            "name": "Encoding",
            "type": [
                "string",
                "null"
            ]
        },
        {
            "name": "EncryptionMethod",
            "type": [
                "EncryptionMethodType",
                "null"
            ]
        },
        {
            "name": "KeyInfo",
            "type": [
                "KeyInfoType",
                "null"
            ]
        },
        {
            "name": "CipherData",
            "type": "CipherDataType"
        },
        {
            "name": "EncryptionProperties",
            "type": [
                "EncryptionPropertiesType",
                "null"
            ]
        },
        {
            "name": "ReferenceList",
            "type": [
                {
                    "type": "record",
                    "name": "type0",
                    "fields": [
                        {
                            "name": "DataReference",
                            "type": [
                                "ReferenceType",
                                "null"
                            ]
                        },
                        {
                            "name": "KeyReference",
                            "type": [
                                "ReferenceType",
                                "null"
                            ]
                        }
                    ]
                },
                "null"
            ]
        },
        {
            "name": "CarriedKeyName",
            "type": [
                "string",
                "null"
            ]
        }
    ]
}`)

// Generated by codegen. Please do not modify.
var _Type0_schema, _Type0_schema_err = avro.ParseSchema(`{
    "type": "record",
    "name": "type0",
    "fields": [
        {
            "name": "DataReference",
            "type": [
                "ReferenceType",
                "null"
            ]
        },
        {
            "name": "KeyReference",
            "type": [
                "ReferenceType",
                "null"
            ]
        }
    ]
}`)

// Generated by codegen. Please do not modify.
var _SubjectConfirmationType_schema, _SubjectConfirmationType_schema_err = avro.ParseSchema(`{
    "type": "record",
    "name": "SubjectConfirmationType",
    "fields": [
        {
            "name": "Method",
            "type": "string"
        },
        {
            "name": "BaseID",
            "type": [
                "BaseIDAbstractType",
                "null"
            ]
        },
        {
            "name": "NameID",
            "type": [
                "NameIDType",
                "null"
            ]
        },
        {
            "name": "EncryptedID",
            "type": [
                "EncryptedElementType",
                "null"
            ]
        },
        {
            "name": "SubjectConfirmationData",
            "type": [
                {
                    "type": "record",
                    "name": "SubjectConfirmationDataType",
                    "fields": [
                        {
                            "name": "NotBefore",
                            "type": [
                                "string",
                                "null"
                            ]
                        },
                        {
                            "name": "NotOnOrAfter",
                            "type": [
                                "string",
                                "null"
                            ]
                        },
                        {
                            "name": "Recipient",
                            "type": [
                                "string",
                                "null"
                            ]
                        },
                        {
                            "name": "InResponseTo",
                            "type": [
                                "string",
                                "null"
                            ]
                        },
                        {
                            "name": "Address",
                            "type": [
                                "string",
                                "null"
                            ]
                        },
                        {
                            "name": "others",
                            "type": {
                                "type": "map",
                                "values": "string"
                            }
                        }
                    ]
                },
                "null"
            ]
        }
    ]
}`)

// Generated by codegen. Please do not modify.
var _SubjectConfirmationDataType_schema, _SubjectConfirmationDataType_schema_err = avro.ParseSchema(`{
    "type": "record",
    "name": "SubjectConfirmationDataType",
    "fields": [
        {
            "name": "NotBefore",
            "type": [
                "string",
                "null"
            ]
        },
        {
            "name": "NotOnOrAfter",
            "type": [
                "string",
                "null"
            ]
        },
        {
            "name": "Recipient",
            "type": [
                "string",
                "null"
            ]
        },
        {
            "name": "InResponseTo",
            "type": [
                "string",
                "null"
            ]
        },
        {
            "name": "Address",
            "type": [
                "string",
                "null"
            ]
        },
        {
            "name": "others",
            "type": {
                "type": "map",
                "values": "string"
            }
        }
    ]
}`)

// Generated by codegen. Please do not modify.
var _NameIDPolicyType_schema, _NameIDPolicyType_schema_err = avro.ParseSchema(`{
    "type": "record",
    "name": "NameIDPolicyType",
    "fields": [
        {
            "name": "Format",
            "type": [
                "string",
                "null"
            ]
        },
        {
            "name": "SPNameQualifier",
            "type": [
                "string",
                "null"
            ]
        },
        {
            "name": "AllowCreate",
            "type": [
                "boolean",
                "null"
            ]
        }
    ]
}`)

// Generated by codegen. Please do not modify.
var _ConditionsType_schema, _ConditionsType_schema_err = avro.ParseSchema(`{
    "type": "record",
    "name": "ConditionsType",
    "fields": [
        {
            "name": "NotBefore",
            "type": [
                "string",
                "null"
            ]
        },
        {
            "name": "NotOnOrAfter",
            "type": [
                "string",
                "null"
            ]
        },
        {
            "name": "Condition",
            "type": [
                {
                    "type": "record",
                    "name": "ConditionAbstractType"
                },
                "null"
            ]
        },
        {
            "name": "AudienceRestriction",
            "type": [
                {
                    "type": "record",
                    "name": "AudienceRestrictionType",
                    "fields": [
                        {
                            "name": "Audience",
                            "type": {
                                "type": "array",
                                "items": "string"
                            }
                        }
                    ]
                },
                "null"
            ]
        },
        {
            "name": "OneTimeUse",
            "type": [
                {
                    "type": "record",
                    "name": "OneTimeUseType"
                },
                "null"
            ]
        },
        {
            "name": "ProxyRestriction",
            "type": [
                {
                    "type": "record",
                    "name": "ProxyRestrictionType",
                    "fields": [
                        {
                            "name": "Count",
                            "type": [
                                "string",
                                "null"
                            ]
                        },
                        {
                            "name": "Audience",
                            "type": {
                                "type": "array",
                                "items": "string"
                            }
                        }
                    ]
                },
                "null"
            ]
        }
    ]
}`)

// Generated by codegen. Please do not modify.
var _ConditionAbstractType_schema, _ConditionAbstractType_schema_err = avro.ParseSchema(`{
    "type": "record",
    "name": "ConditionAbstractType"
}`)

// Generated by codegen. Please do not modify.
var _AudienceRestrictionType_schema, _AudienceRestrictionType_schema_err = avro.ParseSchema(`{
    "type": "record",
    "name": "AudienceRestrictionType",
    "fields": [
        {
            "name": "Audience",
            "type": {
                "type": "array",
                "items": "string"
            }
        }
    ]
}`)

// Generated by codegen. Please do not modify.
var _OneTimeUseType_schema, _OneTimeUseType_schema_err = avro.ParseSchema(`{
    "type": "record",
    "name": "OneTimeUseType"
}`)

// Generated by codegen. Please do not modify.
var _ProxyRestrictionType_schema, _ProxyRestrictionType_schema_err = avro.ParseSchema(`{
    "type": "record",
    "name": "ProxyRestrictionType",
    "fields": [
        {
            "name": "Count",
            "type": [
                "string",
                "null"
            ]
        },
        {
            "name": "Audience",
            "type": {
                "type": "array",
                "items": "string"
            }
        }
    ]
}`)

// Generated by codegen. Please do not modify.
var _RequestedAuthnContextType_schema, _RequestedAuthnContextType_schema_err = avro.ParseSchema(`{
    "type": "record",
    "name": "RequestedAuthnContextType",
    "fields": [
        {
            "name": "Comparison",
            "type": [
                "string",
                "null"
            ]
        },
        {
            "name": "AuthnContextClassRef",
            "type": {
                "type": "array",
                "items": "string"
            }
        },
        {
            "name": "AuthnContextDeclRef",
            "type": {
                "type": "array",
                "items": "string"
            }
        }
    ]
}`)

// Generated by codegen. Please do not modify.
var _ScopingType_schema, _ScopingType_schema_err = avro.ParseSchema(`{
    "type": "record",
    "name": "ScopingType",
    "fields": [
        {
            "name": "ProxyCount",
            "type": [
                "string",
                "null"
            ]
        },
        {
            "name": "IDPList",
            "type": [
                {
                    "type": "record",
                    "name": "IDPListType",
                    "fields": [
                        {
                            "name": "IDPEntry",
                            "type": {
                                "type": "array",
                                "items": {
                                    "type": "record",
                                    "name": "IDPEntryType",
                                    "fields": [
                                        {
                                            "name": "ProviderID",
                                            "type": "string"
                                        },
                                        {
                                            "name": "Name",
                                            "type": [
                                                "string",
                                                "null"
                                            ]
                                        },
                                        {
                                            "name": "Loc",
                                            "type": [
                                                "string",
                                                "null"
                                            ]
                                        }
                                    ]
                                }
                            }
                        },
                        {
                            "name": "GetComplete",
                            "type": [
                                "string",
                                "null"
                            ]
                        }
                    ]
                },
                "null"
            ]
        },
        {
            "name": "RequesterID",
            "type": {
                "type": "array",
                "items": "string"
            }
        }
    ]
}`)

// Generated by codegen. Please do not modify.
var _IDPListType_schema, _IDPListType_schema_err = avro.ParseSchema(`{
    "type": "record",
    "name": "IDPListType",
    "fields": [
        {
            "name": "IDPEntry",
            "type": {
                "type": "array",
                "items": {
                    "type": "record",
                    "name": "IDPEntryType",
                    "fields": [
                        {
                            "name": "ProviderID",
                            "type": "string"
                        },
                        {
                            "name": "Name",
                            "type": [
                                "string",
                                "null"
                            ]
                        },
                        {
                            "name": "Loc",
                            "type": [
                                "string",
                                "null"
                            ]
                        }
                    ]
                }
            }
        },
        {
            "name": "GetComplete",
            "type": [
                "string",
                "null"
            ]
        }
    ]
}`)

// Generated by codegen. Please do not modify.
var _IDPEntryType_schema, _IDPEntryType_schema_err = avro.ParseSchema(`{
    "type": "record",
    "name": "IDPEntryType",
    "fields": [
        {
            "name": "ProviderID",
            "type": "string"
        },
        {
            "name": "Name",
            "type": [
                "string",
                "null"
            ]
        },
        {
            "name": "Loc",
            "type": [
                "string",
                "null"
            ]
        }
    ]
}`)

// Generated by codegen. Please do not modify.
var _ArtifactResponseType_schema, _ArtifactResponseType_schema_err = avro.ParseSchema(`{
    "type": "record",
    "name": "ArtifactResponseType",
    "fields": [
        {
            "name": "ID",
            "type": "string"
        },
        {
            "name": "InResponseTo",
            "type": [
                "string",
                "null"
            ]
        },
        {
            "name": "Version",
            "type": "string"
        },
        {
            "name": "IssueInstant",
            "type": "string"
        },
        {
            "name": "Destination",
            "type": [
                "string",
                "null"
            ]
        },
        {
            "name": "Consent",
            "type": [
                "string",
                "null"
            ]
        },
        {
            "name": "Issuer",
            "type": [
                {
                    "type": "record",
                    "name": "NameIDType",
                    "fields": [
                        {
                            "name": "NameQualifier",
                            "type": [
                                "string",
                                "null"
                            ]
                        },
                        {
                            "name": "SPNameQualifier",
                            "type": [
                                "string",
                                "null"
                            ]
                        },
                        {
                            "name": "Format",
                            "type": [
                                "string",
                                "null"
                            ]
                        },
                        {
                            "name": "SPProvidedID",
                            "type": [
                                "string",
                                "null"
                            ]
                        }
                    ]
                },
                "null"
            ]
        },
        {
            "name": "Signature",
            "type": [
                {
                    "type": "record",
                    "name": "SignatureType",
                    "fields": [
                        {
                            "name": "Id",
                            "type": [
                                "string",
                                "null"
                            ]
                        },
                        {
                            "name": "SignedInfo",
                            "type": {
                                "type": "record",
                                "name": "SignedInfoType",
                                "fields": [
                                    {
                                        "name": "Id",
                                        "type": [
                                            "string",
                                            "null"
                                        ]
                                    },
                                    {
                                        "name": "CanonicalizationMethod",
                                        "type": {
                                            "type": "record",
                                            "name": "CanonicalizationMethodType",
                                            "fields": [
                                                {
                                                    "name": "Algorithm",
                                                    "type": "string"
                                                },
                                                {
                                                    "name": "others",
                                                    "type": {
                                                        "type": "map",
                                                        "values": "string"
                                                    }
                                                }
                                            ]
                                        }
                                    },
                                    {
                                        "name": "SignatureMethod",
                                        "type": {
                                            "type": "record",
                                            "name": "SignatureMethodType",
                                            "fields": [
                                                {
                                                    "name": "Algorithm",
                                                    "type": "string"
                                                },
                                                {
                                                    "name": "HMACOutputLength",
                                                    "type": [
                                                        "string",
                                                        "null"
                                                    ]
                                                },
                                                {
                                                    "name": "others",
                                                    "type": {
                                                        "type": "map",
                                                        "values": "string"
                                                    }
                                                }
                                            ]
                                        }
                                    },
                                    {
                                        "name": "Reference",
                                        "type": {
                                            "type": "array",
                                            "items": {
                                                "type": "record",
                                                "name": "ReferenceType",
                                                "fields": [
                                                    {
                                                        "name": "Id",
                                                        "type": [
                                                            "string",
                                                            "null"
                                                        ]
                                                    },
                                                    {
                                                        "name": "URI",
                                                        "type": [
                                                            "string",
                                                            "null"
                                                        ]
                                                    },
                                                    {
                                                        "name": "Type",
                                                        "type": [
                                                            "string",
                                                            "null"
                                                        ]
                                                    },
                                                    {
                                                        "name": "Transforms",
                                                        "type": [
                                                            {
                                                                "type": "record",
                                                                "name": "TransformsType",
                                                                "fields": [
                                                                    {
                                                                        "name": "Transform",
                                                                        "type": {
                                                                            "type": "array",
                                                                            "items": {
                                                                                "type": "record",
                                                                                "name": "TransformType",
                                                                                "fields": [
                                                                                    {
                                                                                        "name": "Algorithm",
                                                                                        "type": "string"
                                                                                    },
                                                                                    {
                                                                                        "name": "others",
                                                                                        "type": {
                                                                                            "type": "map",
                                                                                            "values": "string"
                                                                                        }
                                                                                    },
                                                                                    {
                                                                                        "name": "XPath",
                                                                                        "type": [
                                                                                            "string",
                                                                                            "null"
                                                                                        ]
                                                                                    }
                                                                                ]
                                                                            }
                                                                        }
                                                                    }
                                                                ]
                                                            },
                                                            "null"
                                                        ]
                                                    },
                                                    {
                                                        "name": "DigestMethod",
                                                        "type": {
                                                            "type": "record",
                                                            "name": "DigestMethodType",
                                                            "fields": [
                                                                {
                                                                    "name": "Algorithm",
                                                                    "type": "string"
                                                                },
                                                                {
                                                                    "name": "others",
                                                                    "type": {
                                                                        "type": "map",
                                                                        "values": "string"
                                                                    }
                                                                }
                                                            ]
                                                        }
                                                    },
                                                    {
                                                        "name": "DigestValue",
                                                        "type": "string"
                                                    }
                                                ]
                                            }
                                        }
                                    }
                                ]
                            }
                        },
                        {
                            "name": "SignatureValue",
                            "type": {
                                "type": "record",
                                "name": "SignatureValueType",
                                "fields": [
                                    {
                                        "name": "Id",
                                        "type": [
                                            "string",
                                            "null"
                                        ]
                                    }
                                ]
                            }
                        },
                        {
                            "name": "KeyInfo",
                            "type": [
                                {
                                    "type": "record",
                                    "name": "KeyInfoType",
                                    "fields": [
                                        {
                                            "name": "Id",
                                            "type": [
                                                "string",
                                                "null"
                                            ]
                                        },
                                        {
                                            "name": "KeyName",
                                            "type": [
                                                "string",
                                                "null"
                                            ]
                                        },
                                        {
                                            "name": "KeyValue",
                                            "type": [
                                                {
                                                    "type": "record",
                                                    "name": "KeyValueType",
                                                    "fields": [
                                                        {
                                                            "name": "DSAKeyValue",
                                                            "type": [
                                                                {
                                                                    "type": "record",
                                                                    "name": "DSAKeyValueType",
                                                                    "fields": [
                                                                        {
                                                                            "name": "P",
                                                                            "type": "string"
                                                                        },
                                                                        {
                                                                            "name": "Q",
                                                                            "type": "string"
                                                                        },
                                                                        {
                                                                            "name": "G",
                                                                            "type": [
                                                                                "string",
                                                                                "null"
                                                                            ]
                                                                        },
                                                                        {
                                                                            "name": "Y",
                                                                            "type": "string"
                                                                        },
                                                                        {
                                                                            "name": "J",
                                                                            "type": [
                                                                                "string",
                                                                                "null"
                                                                            ]
                                                                        },
                                                                        {
                                                                            "name": "Seed",
                                                                            "type": "string"
                                                                        },
                                                                        {
                                                                            "name": "PgenCounter",
                                                                            "type": "string"
                                                                        }
                                                                    ]
                                                                },
                                                                "null"
                                                            ]
                                                        },
                                                        {
                                                            "name": "RSAKeyValue",
                                                            "type": [
                                                                {
                                                                    "type": "record",
                                                                    "name": "RSAKeyValueType",
                                                                    "fields": [
                                                                        {
                                                                            "name": "Modulus",
                                                                            "type": "string"
                                                                        },
                                                                        {
                                                                            "name": "Exponent",
                                                                            "type": "string"
                                                                        }
                                                                    ]
                                                                },
                                                                "null"
                                                            ]
                                                        },
                                                        {
                                                            "name": "others",
                                                            "type": {
                                                                "type": "map",
                                                                "values": "string"
                                                            }
                                                        }
                                                    ]
                                                },
                                                "null"
                                            ]
                                        },
                                        {
                                            "name": "RetrievalMethod",
                                            "type": [
                                                {
                                                    "type": "record",
                                                    "name": "RetrievalMethodType",
                                                    "fields": [
                                                        {
                                                            "name": "URI",
                                                            "type": [
                                                                "string",
                                                                "null"
                                                            ]
                                                        },
                                                        {
                                                            "name": "Type",
                                                            "type": [
                                                                "string",
                                                                "null"
                                                            ]
                                                        },
                                                        {
                                                            "name": "Transforms",
                                                            "type": [
                                                                "TransformsType",
                                                                "null"
                                                            ]
                                                        }
                                                    ]
                                                },
                                                "null"
                                            ]
                                        },
                                        {
                                            "name": "X509Data",
                                            "type": [
                                                {
                                                    "type": "record",
                                                    "name": "X509DataType",
                                                    "fields": [
                                                        {
                                                            "name": "X509IssuerSerial",
                                                            "type": [
                                                                {
                                                                    "type": "record",
                                                                    "name": "X509IssuerSerialType",
                                                                    "fields": [
                                                                        {
                                                                            "name": "X509IssuerName",
                                                                            "type": "string"
                                                                        },
                                                                        {
                                                                            "name": "X509SerialNumber",
                                                                            "type": "string"
                                                                        }
                                                                    ]
                                                                },
                                                                "null"
                                                            ]
                                                        },
                                                        {
                                                            "name": "X509SKI",
                                                            "type": [
                                                                "string",
                                                                "null"
                                                            ]
                                                        },
                                                        {
                                                            "name": "X509SubjectName",
                                                            "type": [
                                                                "string",
                                                                "null"
                                                            ]
                                                        },
                                                        {
                                                            "name": "X509Certificate",
                                                            "type": [
                                                                "string",
                                                                "null"
                                                            ]
                                                        },
                                                        {
                                                            "name": "X509CRL",
                                                            "type": [
                                                                "string",
                                                                "null"
                                                            ]
                                                        },
                                                        {
                                                            "name": "others",
                                                            "type": {
                                                                "type": "map",
                                                                "values": "string"
                                                            }
                                                        }
                                                    ]
                                                },
                                                "null"
                                            ]
                                        },
                                        {
                                            "name": "PGPData",
                                            "type": [
                                                {
                                                    "type": "record",
                                                    "name": "PGPDataType",
                                                    "fields": [
                                                        {
                                                            "name": "PGPKeyID",
                                                            "type": [
                                                                "string",
                                                                "null"
                                                            ]
                                                        },
                                                        {
                                                            "name": "PGPKeyPacket0",
                                                            "type": [
                                                                "string",
                                                                "null"
                                                            ]
                                                        },
                                                        {
                                                            "name": "others",
                                                            "type": {
                                                                "type": "map",
                                                                "values": "string"
                                                            }
                                                        }
                                                    ]
                                                },
                                                "null"
                                            ]
                                        },
                                        {
                                            "name": "SPKIData",
                                            "type": [
                                                {
                                                    "type": "record",
                                                    "name": "SPKIDataType",
                                                    "fields": [
                                                        {
                                                            "name": "SPKISexp",
                                                            "type": "string"
                                                        },
                                                        {
                                                            "name": "others",
                                                            "type": {
                                                                "type": "map",
                                                                "values": "string"
                                                            }
                                                        }
                                                    ]
                                                },
                                                "null"
                                            ]
                                        },
                                        {
                                            "name": "MgmtData",
                                            "type": [
                                                "string",
                                                "null"
                                            ]
                                        },
                                        {
                                            "name": "others",
                                            "type": {
                                                "type": "map",
                                                "values": "string"
                                            }
                                        }
                                    ]
                                },
                                "null"
                            ]
                        },
                        {
                            "name": "Object",
                            "type": {
                                "type": "array",
                                "items": {
                                    "type": "record",
                                    "name": "ObjectType",
                                    "fields": [
                                        {
                                            "name": "Id",
                                            "type": [
                                                "string",
                                                "null"
                                            ]
                                        },
                                        {
                                            "name": "MimeType",
                                            "type": [
                                                "string",
                                                "null"
                                            ]
                                        },
                                        {
                                            "name": "Encoding",
                                            "type": [
                                                "string",
                                                "null"
                                            ]
                                        },
                                        {
                                            "name": "others",
                                            "type": {
                                                "type": "map",
                                                "values": "string"
                                            }
                                        }
                                    ]
                                }
                            }
                        }
                    ]
                },
                "null"
            ]
        },
        {
            "name": "Extensions",
            "type": [
                {
                    "type": "record",
                    "name": "ExtensionsType",
                    "fields": [
                        {
                            "name": "others",
                            "type": {
                                "type": "map",
                                "values": "string"
                            }
                        }
                    ]
                },
                "null"
            ]
        },
        {
            "name": "Status",
            "type": {
                "type": "record",
                "name": "StatusType",
                "fields": [
                    {
                        "name": "StatusCode",
                        "type": {
                            "type": "record",
                            "name": "StatusCodeType",
                            "fields": [
                                {
                                    "name": "Value",
                                    "type": "string"
                                },
                                {
                                    "name": "StatusCode",
                                    "type": [
                                        "StatusCodeType",
                                        "null"
                                    ]
                                }
                            ]
                        }
                    },
                    {
                        "name": "StatusMessage",
                        "type": [
                            "string",
                            "null"
                        ]
                    },
                    {
                        "name": "StatusDetail",
                        "type": [
                            {
                                "type": "record",
                                "name": "StatusDetailType",
                                "fields": [
                                    {
                                        "name": "others",
                                        "type": {
                                            "type": "map",
                                            "values": "string"
                                        }
                                    }
                                ]
                            },
                            "null"
                        ]
                    }
                ]
            }
        },
        {
            "name": "others",
            "type": {
                "type": "map",
                "values": "string"
            }
        }
    ]
}`)

// Generated by codegen. Please do not modify.
var _StatusType_schema, _StatusType_schema_err = avro.ParseSchema(`{
    "type": "record",
    "name": "StatusType",
    "fields": [
        {
            "name": "StatusCode",
            "type": {
                "type": "record",
                "name": "StatusCodeType",
                "fields": [
                    {
                        "name": "Value",
                        "type": "string"
                    },
                    {
                        "name": "StatusCode",
                        "type": [
                            "StatusCodeType",
                            "null"
                        ]
                    }
                ]
            }
        },
        {
            "name": "StatusMessage",
            "type": [
                "string",
                "null"
            ]
        },
        {
            "name": "StatusDetail",
            "type": [
                {
                    "type": "record",
                    "name": "StatusDetailType",
                    "fields": [
                        {
                            "name": "others",
                            "type": {
                                "type": "map",
                                "values": "string"
                            }
                        }
                    ]
                },
                "null"
            ]
        }
    ]
}`)

// Generated by codegen. Please do not modify.
var _StatusCodeType_schema, _StatusCodeType_schema_err = avro.ParseSchema(`{
    "type": "record",
    "name": "StatusCodeType",
    "fields": [
        {
            "name": "Value",
            "type": "string"
        },
        {
            "name": "StatusCode",
            "type": [
                "StatusCodeType",
                "null"
            ]
        }
    ]
}`)

// Generated by codegen. Please do not modify.
var _StatusDetailType_schema, _StatusDetailType_schema_err = avro.ParseSchema(`{
    "type": "record",
    "name": "StatusDetailType",
    "fields": [
        {
            "name": "others",
            "type": {
                "type": "map",
                "values": "string"
            }
        }
    ]
}`)
