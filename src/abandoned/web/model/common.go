package model

var (
	//Cl, Errr              = client.Dial(client.Options{})
	ServiceServer *Server = nil
)

type Decimate struct { //Reduce the polygon count by removing some vertices
	//The maximum distance between the surface vertices and the simplified surface
	SurfacicTolerance float64 `json:"surfacicTolerance,omitempty"`
	//The maximum distance between the vertex of the straight line and the simplified straight line
	LineicTolerance float64 `json:"lineicTolerance,omitempty"`
	//The maximum angle between the original normal and the normal interpolated on the simplified surface
	NormalTolerance float64 `json:"normalTolerance,omitempty"`
	//In UV space, the maximum distance between the original texcord and those interpolated on the simplified surface
	TexCoordTolerance float64 `json:"texCoordTolerance,omitempty"`
	//if for True to not process normals and/or in very small areas texcoord toleranceConstraintsAccordingTo surfacicTolerance）
	ReleaseConstraintOnSmallArea float64 `json:"releaseConstraintOnSmallArea,omitempty"`
}

// Repair CAD shapes, assemble faces, remove duplicate faces, optimize loops and repair topology
type RepairCAD struct {
	Tolerance float64 `json:"tolerance,omitempty"`
	Orient    bool    `json:"orient,omitempty"`
}

// cadDoesATessellationForEachGivenPart
type Tessellate struct {
	//Maximum distance between geometry and subdivision surface
	MaxSag float64 `json:"maxSag,omitempty"`
	//maximumLengthOfElement
	MaxLength int `json:"maxLength,omitempty"`
	//The maximum angle between the normals of two adjacent elements
	MaxAngle int `json:"maxAngle,omitempty"`
	//If true, generate normals
	CreateNormals bool `json:"createNormals,omitempty"`
	//Select texture coordinate generation mode NoUV (0)  FastUV (1)  UniformUV (2)
	UvMode int `json:"uvMode,omitempty"`
	//Generate UV channel of texture coordinates
	UvChannel int `json:"uvChannel,omitempty"`
	//UV filling between UV islands in UV coordinate space (between 0-1).
	UvPadding int `json:"uvPadding,omitempty"`
	//If true, tangents will be generated
	CreateTangents bool `json:"createTangents,omitempty"`
	//If true, free edges will be created for each Patch boundary
	CreateFreeEdges bool `json:"createFreeEdges,omitempty" default:"true"`
	//If true, the BRep shape will be retained for Back to BRep or Retessellate
	KeepBRepShape bool `json:"keepBRepShape,omitempty"`
	//If true, the parts that have been subdivided will be re-subdivided.
	OverrideExistingTessellation bool `json:"overrideExistingTessellation,omitempty"`
}

type RepairMesh struct {
	// connectionTolerance
	Tolerance float64 `json:"tolerance,omitempty"`
	//At the end of the repair process, the crack results in a non-manifold edge
	CrackNonManifold bool `json:"crackNonManifold,omitempty"`
	//ifTrueRepositionTheModel
	Orient bool `json:"orient,omitempty"`
}

type OptimizeFunc struct {
	RepairMesh RepairMesh `json:"repairMesh,omitempty"`
	Decimate   Decimate   `json:"decimate,omitempty"`
	RepairCAD  RepairCAD  `json:"repairCAD,omitempty"`
	Tessellate Tessellate `json:"tessellate,omitempty"`
}

type NodeMerges struct {
	Name            string `json:"name,omitempty"`            //The name of the node that needs to be merged
	CaseSensitive   bool   `json:"caseSensitive,omitempty"`   //Case Sensitive
	MatchWholeWorld bool   `json:"matchWholeWorld,omitempty"` //Full word match
	RegExp          bool   `json:"regExp,omitempty"`          //Enable regular expressions
}

type NodeDecimate struct {
	TargetStrategy    []interface{} `json:"TargetStrategy,omitempty"`    /// Face reduction strategy: specify the number of faces ["triangle Count":10000], proportionally ["ratio",50.000000]
	BoundaryWeight    float32       `json:"boundaryWeight,omitempty"`    /// 1.000000,//Boundary weight
	NormalWeight      float32       `json:"normalWeight,omitempty"`      /// 1.000000,//Normal vector weight
	UVWeight          float32       `json:"UVWeight,omitempty"`          /// 1.000000,//uv weights
	SharpNormalWeight float32       `json:"sharpNormalWeight,omitempty"` /// 1.000000,
	UVSeamWeight      float32       `json:"UVSeamWeight,omitempty"`      /// 10.000000,//uv seamWeight
	ForbidUVFoldovers bool          `json:"forbidUVFoldovers,omitempty"` /// True, // prohibit UV overlap
	ProtectTopology   bool          `json:"protectTopology,omitempty"`   /// True//Guarantee topology
}

type NodeDecimates struct {
	NodeMerge    NodeMerges   `json:"nodeMerge,omitempty"`
	NodeDecimate NodeDecimate `json:"nodeDecimate,omitempty"`
}

type ProcessMesh struct {
	NodeMerges    []NodeMerges    `json:"nodeMerges,omitempty"`
	NodeDecimates []NodeDecimates `json:"nodeDecimates,omitempty"`
}
