package model_function

type ModelFrameLayerType string

const ModelFrameLayerTypeIO = ModelFrameLayerType("io")

const ModelFrameLayerTypeLogicRepo = ModelFrameLayerType("logic_repo")
const ModelFrameLayerTypeDataRepo = ModelFrameLayerType("data_repo")

const ModelFrameLayerTypeLogicRelay = ModelFrameLayerType("logic_relay")
const ModelFrameLayerTypeDataRelay = ModelFrameLayerType("data_relay")

const ModelFrameLayerTypeLogicClient = ModelFrameLayerType("logic_client")
const ModelFrameLayerTypeDataClient = ModelFrameLayerType("data_client")
