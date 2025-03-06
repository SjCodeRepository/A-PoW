package config

const (
	EndCodeNum              = 1000 // Number of experiment rounds
	TotalJudgeNum           = 10   //Total count of adaptive mining difficulty adjustments
	TotalNodeNum            = 30   //Total number of nodes in the network
	Nnum                    = 2    //Total number of mining difficulties
	Nmin                    = 8    //Minimum mining difficulty value
	Mmax                    = 9    //Maximum mining difficulty value
	MaxSNum                 = 10   // Maximum successful mining count
	Alpha                   = 0.1  // Learning rate
	Gamma                   = 0.9  // Discount factor
	Epsilon                 = 0.1  // Exploration factor
	Epochs                  = 5    // Number of iterations
	Range                   = 2    //The difficulty selection range of mining difficulty adjustment algorithm
	NodeNum                 = 5    //The default number of nodes in each node group
	NodeGroups              = 6    //Total number of node groups
	ActionsNumber           = 20   //Total number of actions
	Max_Epoch               = 2000 //The maximum number of iterations in the mining difficulty initialization adjustment phase
	Max_Variance            = 3000 //Maximum reasonable variance
	Beta                    = 0.9
	R1                  int = 5  //reward 1
	R2                  int = -1 //reward 2
	Range1                  = 1.0
	Range2                  = 0.8
	Max_Range               = 0.5
	CRC                     = 100 //The total number of rounds between two mining difficulty adjustments
	RandP                   = 0.8
	CalculateTVRoundNum     = 5.0 //Calculate the total number of rounds of reference trust value
)

var Target []byte = []byte{'0', '1'} //List of mining targets
var NodeTimes [][]float64            //Average mining time of each node group
var TestReferenceTrustValue = []float64{0.9318024793453608, 0.8724570355615446, 0.9657830669857238, 0.9878232994904835, 0.8574912084146933, 0.9594154126038474, 0.9118622353898499, 0.947111004593656, 0.9089607995833419, 0.9562192432413789, 0.9259529051523792, 0.8925467477592449, 0.46544912596244825, 0.5031891346287077, 0.4315250183263282, 0.4658121517978799, 0.5084992182060362, 0.5236693640908544, 0.4133962029076068, 0.5225435088124969, 0.5569367014706129, 0.4637775450739829, 0.45726215521488095, 0.5944979546140724, 0.03000738195379809, 0.11430630579810167, 0.017925662295422275, 0.04395332545483821, 0.13423839627424275, 0.06670415145578028}

// network layer
var (
	Delay       int // The delay of network (ms) when sending. 0 if delay < 0
	JitterRange int // The jitter range of delay (ms). Jitter follows a uniform distribution. 0 if JitterRange < 0.
	Bandwidth   int // The bandwidth limit (Bytes). +inf if bandwidth < 0
)
