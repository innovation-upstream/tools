package model_frame_path

type ModelFramePathType string

const ModelFramePathTypeCreate = ModelFramePathType("create")
const ModelFramePathTypeUpdate = ModelFramePathType("update")
const ModelFramePathTypeReadOne = ModelFramePathType("read_one")
const ModelFramePathTypeReadPage = ModelFramePathType("read_page")
const ModelFramePathTypeReadAll = ModelFramePathType("read_all")
const ModelFramePathTypeDelete = ModelFramePathType("delete")
