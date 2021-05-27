package model_frame_path

type ModelFunctionType string

const ModelFunctionTypeCreate = ModelFunctionType("create")
const ModelFunctionTypeUpdate = ModelFunctionType("update")
const ModelFunctionTypeReadOne = ModelFunctionType("read_one")
const ModelFunctionTypeReadPage = ModelFunctionType("read_page")
const ModelFunctionTypeReadAll = ModelFunctionType("read_all")
const ModelFunctionTypeDelete = ModelFunctionType("delete")
