package model_frame_path

type ReferenceType string

// Reference the model owner
const OwnerReferenceType = ReferenceType("owner")

// Reference an external model
const ExternalReferenceType = ReferenceType("external")

// Reference the model (self)
const InternalReferenceType = ReferenceType("internal")
