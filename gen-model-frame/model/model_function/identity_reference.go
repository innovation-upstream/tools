package model_function

type ReferenceType string

// Reference the model owner
const OwnerReferenceType = ReferenceType("owner")

// Reference an external model
const ExternalReferenceType = ReferenceType("external")

// Reference the model (self)
const InternalReferenceType = ReferenceType("internal")
