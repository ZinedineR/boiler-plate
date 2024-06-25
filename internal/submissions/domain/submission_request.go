package domain

const (
	ProfileRiskCategoryConservative = "Conservative"
	ProfileRiskCategoryModerate     = "Moderate"
	ProfileRiskCategoryBalanced     = "Balanced"
	ProfileRiskCategoryGrowth       = "Growth"
	ProfileRiskCategoryAggressive   = "Aggressive"
)

type SubmissionRequest struct {
	UserId  int `validate:"required" json:"user_id"`
	Answers []struct {
		QuestionId int    `json:"question_id"`
		Answer     string `json:"answer"`
	} `json:"answers"`
}

type SubmissionResult struct {
	UserId     int    `json:"user_id"`
	QuestionId int    `json:"question_id"`
	Answer     string `json:"answer"`
}

func (model *Submissions) DeclareRisk(answers []SubmissionResult) {
	score := 0
	for _, answer := range answers {
		for _, question := range Questions {
			if question.ID == answer.QuestionId {
				for _, option := range question.Options {
					if answer.Answer == option.Answer {
						score += option.Weight
					}
				}
			}
		}
	}
	model.RiskScore = score
	for _, profile := range RiskMapping {
		if score >= profile.MinScore && score <= profile.MaxScore {
			model.RiskCategory = profile.Category
			model.RiskDefinition = profile.Definition
		}
	}
}

var Questions = []Question{
	{
		ID:       1,
		Question: "Apakah tujuan investasi Anda?",
		Options: []Option{
			{Answer: "Pertumbuhan kekayaan untuk jangka panjang", Weight: 5},
			{Answer: "Pendapatan dan pertumbuhan dalam jangka panjang", Weight: 4},
			{Answer: "Pendapatan berkala", Weight: 3},
			{Answer: "Pendapatan dan keamanan dana investasi", Weight: 2},
			{Answer: "Keamanan dana investasi", Weight: 1},
		},
	},
	{
		ID:       2,
		Question: "Berdasarkan tujuan investasi Anda, dana Anda akan diinvestasikan untuk jangka waktu?",
		Options: []Option{
			{Answer: "≥ 10 tahun", Weight: 5},
			{Answer: "7 - 10 tahun", Weight: 4},
			{Answer: "4 - ≥ 6 tahun", Weight: 3},
			{Answer: "1 - ≥ 3 tahun", Weight: 2},
			{Answer: "< 1 tahun", Weight: 1},
		},
	},
	{
		ID:       3,
		Question: "Berapa lama pengalaman Anda berinvestasi dalam produk yang nilainya berfluktuasi?",
		Options: []Option{
			{Answer: "> 10 tahun", Weight: 5},
			{Answer: "8 - 10 tahun", Weight: 4},
			{Answer: "4 - 7 tahun", Weight: 3},
			{Answer: "< 4 tahun", Weight: 2},
			{Answer: "0 tahun (tidak memiliki pengalaman)", Weight: 1},
		},
	},
	{
		ID:       4,
		Question: "Jenis investasi apa yang pernah Anda miliki?",
		Options: []Option{
			{Answer: "Saham, Reksa Dana terbuka, equity linked structure product", Weight: 5},
			{Answer: "Mata uang asing, currency linked structured product", Weight: 4},
			{Answer: "Uang tunai, deposito, produk dengan proteksi modal", Weight: 3},
		},
	},
	{
		ID:       5,
		Question: "Berapa persen dari aset Anda yang disimpan dalam produk investasi berfluktuasi?",
		Options: []Option{
			{Answer: "> 50%", Weight: 5},
			{Answer: "> 25% - ≥ 50%", Weight: 4},
			{Answer: "> 10% - ≥ 25%", Weight: 3},
			{Answer: "> 0% - ≥ 10%", Weight: 2},
			{Answer: "0%", Weight: 1},
		},
	},
	{
		ID:       6,
		Question: "Tingkat kenaikan dan penurunan nilai investasi yang dapat Anda terima?",
		Options: []Option{
			{Answer: "< -20% - > +20%", Weight: 5},
			{Answer: "-20% - +20%", Weight: 4},
			{Answer: "-15% - +15%", Weight: 3},
			{Answer: "-10% - +10%", Weight: 2},
			{Answer: "-5% - +5%", Weight: 1},
		},
	},
	{
		ID:       7,
		Question: "Ketergantungan Anda pada hasil investasi untuk biaya hidup sehari-hari?",
		Options: []Option{
			{Answer: "Tidak bergantung pada hasil investasi", Weight: 5},
			{Answer: "Tidak bergantung pada hasil investasi, minimal 5 tahun ke depan", Weight: 4},
			{Answer: "Sedikit bergantung pada hasil investasi", Weight: 3},
			{Answer: "Bergantung pada hasil investasi", Weight: 2},
			{Answer: "Sangat bergantung pada hasil investasi", Weight: 1},
		},
	},
	{
		ID:       8,
		Question: "Persentase pendapatan bulanan yang dapat Anda sisihkan untuk investasi/tabungan?",
		Options: []Option{
			{Answer: "> 50%", Weight: 5},
			{Answer: "> 25% - 50%", Weight: 4},
			{Answer: "> 10% - 25%", Weight: 3},
			{Answer: "> 0% - 10%", Weight: 2},
			{Answer: "0%", Weight: 1},
		},
	},
}

var RiskMapping = []RiskProfile{
	{
		MinScore: 0,
		MaxScore: 11,
		Category: ProfileRiskCategoryConservative,
		Definition: "Tujuan utama Anda adalah untuk melindungi modal/dana yang ditempatkan dan Anda tidak memiliki toleransi " +
			"sama sekali terhadap perubahan harga/nilai dari dana investasinya tersebut. " +
			"Anda memiliki pengalaman yang sangat terbatas atau tidak memiliki pengalaman sama sekali mengenai produk investasi.",
	},
	{
		MinScore:   12,
		MaxScore:   19,
		Category:   ProfileRiskCategoryModerate,
		Definition: "Anda memiliki toleransi yang rendah dengan perubahan harga/nilai dari dana investasi dan risiko investasi.",
	},
	{
		MinScore: 20,
		MaxScore: 28,
		Category: ProfileRiskCategoryBalanced,
		Definition: "Anda memiliki toleransi yang cukup terhadap produk investasi dan dapat menerima perubahan yang besar dari " +
			"harga/nilai dari harga yang diinvestasikan.",
	},
	{
		MinScore: 29,
		MaxScore: 35,
		Category: ProfileRiskCategoryGrowth,
		Definition: "Anda memiliki toleransi yang cukup tinggi dan dapat menerima perubahan yang besar dari harga/nilai portfolio" +
			"pada produk investasi yang diinvestasikan." +
			"Pada umumnya Anda sudah pernah atau berpengalaman dalam berinvestasi di produk investasi.",
	},
	{
		MinScore: 36,
		MaxScore: 40,
		Category: ProfileRiskCategoryAggressive,
		Definition: "Anda sangat berpengalaman terhadap produk investasi dan memiliki toleransi yang sangat tinggi atas" +
			"produk-produk investasi. Anda bahkan dapat menerima perubahan signifikan pada modal/nilai investasi." +
			"Pada umumnya portfolio Anda sebagian besar dialokasikan pada produk investasi.",
	},
}
