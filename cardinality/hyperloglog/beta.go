package hyperloglog

import "math"

// beta-function correction coefficients
var beta = [15][8]uint64{
	{0xbfe2a481c6ff465f, 0xbffef6fd81fa8535, 0x4026289d1fa51ce8, 0xc03621a0a440907f, 0x403681615c2cf0e9, 0xc028005edfde99d4, 0x4009c3bf2302eb7c, 0xbfd5e704f4479845},
	{0xbfe80f9079c7fb1b, 0xbfeeb02713651449, 0x401666217da43a02, 0xc0206b6625ce9312, 0x401a0958313af164, 0xc00576d81aa8cb99, 0x3fe1f6149516335c, 0xbfa7b8faf0e84898},
	{0x403dd366fad3cff3, 0xc03f54263ab779d7, 0xc0253041d68391cb, 0xc02724ded4d8bd0c, 0x400e8d0e90b8ba4b, 0xc00353fea4168293, 0x3fdd11f48353a36b, 0xbfad72ad734f0ae9},
	{0x40067b7a7094e84c, 0xc00fd30bccfb505b, 0x3ff50f6f09e881d5, 0xc00f66e8c16f2db7, 0x4000108e1f1fa359, 0xbfe8163e2c8c41aa, 0x3fc03304f9542849, 0xbf86845c1ee9a660},
	{0x3ff019f3331b9cc6, 0xc0000be45d419666, 0x3ffa4c95be3cf174, 0xc005a515dbf5eb97, 0x3ff6460a71441e66, 0xbfddbdb4c4546b94, 0x3fb2e75d0125e52d, 0xbf77b29671c6fafd},
	{0xbfb81aa530881d68, 0xbfe9007d5300299b, 0x3ffb7140916b1e5f, 0xbffbcb367b26ebb3, 0x3feba949d2f8c246, 0xbfce7d04d7ecad35, 0x3fa11e531eb34c5d, 0xbf61071bf9d4928d},
	{0xbfd099418c396868, 0xbfe0d4da57a4933d, 0x3ff7d44c0f701b50, 0xbff4be2a62b98e60, 0x3fe3ee5e02bd892d, 0xbfc40f820ea6fb03, 0x3f9509872cc818cb, 0xbf526e1be62ade83},
	{0xbfdbab38ccc25986, 0xbfbbc36d6e1043e9, 0x3fe37e35e0f27909, 0xbf90f768041bc4bb, 0xbfb45f8c13511d5f, 0x3fa82860c90e8aba, 0xbf8000a4f4ecc346, 0x3f4325337f5d7e42},
	{0xbfd8a37fcf308623, 0x3fc771dc2c8b54ad, 0x3fc0b0d6b50a8dc2, 0x3fb20b3b7e4fc124, 0xbf82591088a65dec, 0x3f8724fb098a3aeb, 0xbf5fd4ed6b2fcfee, 0x3f2d8c5e8db6d9d5},
	{0xbfdaa8ccb20d7d58, 0xbfcc5905ec062a5f, 0x3fd8df2becca188b, 0x3fdd04aa86f1e114, 0xbfd7359d5ec10a0e, 0x3fbf7ff9b7b00597, 0xbf916c7c6c2df81d, 0x3f50d5a9aa68b9bf},
	{0xbfd7be9fb8ac0206, 0x3f840bcb261d17a6, 0x3fc7c82c45c5ccc7, 0x3fc9fc69ad3a57f9, 0xbfbde0bda23ef507, 0x3fa6129a0b14e5db, 0xbf788f16cfa289a6, 0x3f3d78c9425f19c4},
	{0xbfd8752b60ce9b8e, 0xbfec8090b6f89653, 0x3fd810c446037e01, 0x3fefc99a7007a7e7, 0xbfe4fc1aa0fd6aca, 0x3fc777245187e282, 0xbf96f40912ae48b0, 0x3f53e3e04ff1fcbd},
	{0xbfd7e474653543b0, 0xbff6ac32f277ed22, 0x3fda1111d3bc484b, 0x3ff8fbfcc04539d7, 0xbfefc1ec7c43c2d7, 0x3fd0ae6ff97a856e, 0xbf9f45621e70082f, 0x3f598579e561f1a5},
	{0xbfd7894c5d2cbef8, 0x3fe139dec01ed380, 0x3fe8a167f6051c12, 0x3fe199cfc81fc79a, 0xbfe7dd3b70d3b542, 0x3fd074a08e8cf21d, 0xbfa19a20b1442bcc, 0x3f5e7745c186cade},
	{0xbfd758d24ce3749e, 0x3fefe9ea5635c955, 0x3ff8db50cf17bc7b, 0x3ff42633d44a1f18, 0xbff8883b183dcf95, 0x3fde97b90668ef09, 0xbfae7820b1e8f330, 0x3f67d852b2199895},
}

func betaEstimation(p uint64, z float64) float64 {
	zl := math.Log(z + 1)
	_ = beta[14]
	return math.Float64frombits(beta[p][0])*z + math.Float64frombits(beta[p][1])*zl +
		math.Float64frombits(beta[p][2])*math.Pow(zl, 2) + math.Float64frombits(beta[p][3])*math.Pow(zl, 3) +
		math.Float64frombits(beta[p][4])*math.Pow(zl, 4) + math.Float64frombits(beta[p][5])*math.Pow(zl, 5) +
		math.Float64frombits(beta[p][6])*math.Pow(zl, 6) + math.Float64frombits(beta[p][7])*math.Pow(zl, 7)
}
