package reqvo

type Validator interface {
	IsValid() bool
}

type CreateDBReqVO struct {
	Namespace string `json:"namespace"`
	DbName    string `json:"dbName"`
}

func (r *CreateDBReqVO) IsValid() bool {
	return r.Namespace != "" && r.DbName != ""
}

type ExecuteCommandReqVO struct {
	Namespace string `json:"namespace"`
	DbName    string `json:"dbName"`
	Cmd       string `json:"cmd"`
}

func (r *ExecuteCommandReqVO) IsValid() bool {
	return r.Namespace != "" && r.DbName != "" && r.Cmd != ""
}

type QueryCommandReqVO struct {
	Namespace string `json:"namespace"`
	DbName    string `json:"dbName"`
	Cmd       string `json:"cmd"`
}

func (r *QueryCommandReqVO) IsValid() bool {
	return r.Namespace != "" && r.DbName != "" && r.Cmd != ""
}

type DropDBReqVO struct {
	Namespace string `json:"namespace"`
	DbName    string `json:"dbName"`
}

func (r *DropDBReqVO) IsValid() bool {
	return r.Namespace != "" && r.DbName != ""
}

type NewNamespaceReqVO struct {
	Namespace string `json:"namespace"`
}

func (r *NewNamespaceReqVO) IsValid() bool {
	return r.Namespace != ""
}

type DeleteNamespaceReqVO struct {
	Namespace string `json:"namespace"`
}

func (r *DeleteNamespaceReqVO) IsValid() bool {
	return r.Namespace != ""
}

type ShowNamespaceReqVO struct {
	Namespace string `json:"namespace"`
}

func (r *ShowNamespaceReqVO) IsValid() bool {
	return r.Namespace != ""
}

type GetDBSizeReqVO struct {
	Namespace string `json:"namespace"`
	DbName    string `json:"dbName"`
}

func (r *GetDBSizeReqVO) IsValid() bool {
	return r.Namespace != "" && r.DbName != ""
}
