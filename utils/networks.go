package utils

import (
	"math/big"
)

// Network represents the network as a string.
type Network string

// NetworkID represents the network ID as an unsigned 64-bit integer.
type NetworkID uint64

// Predefined network IDs and network strings.
var (
	// EthereumNetworkID is the ID for the Ethereum network.
	EthereumNetworkID NetworkID = 1
	EthereumNetwork Network = "ethereum"

	// GoerliNetworkID is the ID for the Goerli network.
	GoerliNetworkID NetworkID = 5
	GoerliNetwork Network = "goerli"

	// OptimismNetworkID is the ID for the Optimism network.
	OptimismNetworkID NetworkID = 10
	OptimismNetwork Network = "optimism"

	// FlareNetworkID is the ID for the Flare network.
	FlareNetworkID NetworkID = 14
	FlareNetwork Network = "flare"

	// RootstockNetworkID is the ID for the Rootstock network.
	RootstockNetworkID NetworkID = 30
	RootstockNetwork Network = "rootstock"

	// LuksoNetworkID is the ID for the Lukso network.
	LuksoNetworkID NetworkID = 42
	LuksoNetwork Network = "lukso"

	// CrabNetworkID is the ID for the Crab network.
	CrabNetworkID NetworkID = 44
	CrabNetwork Network = "crab"

	// DarwiniaNetworkID is the ID for the Darwinia network.
	DarwiniaNetworkID NetworkID = 46
	DarwiniaNetwork Network = "darwinia"

	// BscNetworkID is the ID for the BSC network.
	BscNetworkID NetworkID = 56
	BscNetwork Network = "bsc"

	// GnosisNetworkID is the ID for the Gnosis network.
	GnosisNetworkID NetworkID = 100
	GnosisNetwork Network = "gnosis"

	// PolygonNetworkID is the ID for the Polygon network.
	PolygonNetworkID NetworkID = 137
	PolygonNetwork Network = "polygon"

	// ShimmerEVMNetworkID is the ID for the Shimmer EVM network.
	ShimmerEVMNetworkID NetworkID = 148
	ShimmerEVMNetwork Network = "shimmerevm"

	// MantaNetworkID is the ID for the Manta network.
	MantaNetworkID NetworkID = 169
	MantaNetwork Network = "manta"

	// XLayerTestnetNetworkID is the ID for the XLayer Testnet network.
	XLayerTestnetNetworkID NetworkID = 195
	XLayerTestnetNetwork Network = "xlayertestnet"

	// XLayerNetworkID is the ID for the XLayer network.
	XLayerNetworkID NetworkID = 196
	XLayerNetwork Network = "xlayer"

	// FantomNetworkID is the ID for the Fantom network.
	FantomNetworkID NetworkID = 250
	FantomNetwork Network = "fantom"

	// KromaNetworkID is the ID for the Kroma network.
	KromaNetworkID NetworkID = 255
	KromaNetwork Network = "kroma"

	// BobaNetworkID is the ID for the Boba network.
	BobaNetworkID NetworkID = 288
	BobaNetwork Network = "boba"

	// ZksyncEraNetworkID is the ID for the Zksync Era network.
	ZksyncEraNetworkID NetworkID = 324
	ZksyncEraNetwork Network = "zksyncera"

	// PublicGoodsNetworkID is the ID for the Public Goods network.
	PublicGoodsNetworkID NetworkID = 424
	PublicGoodsNetwork Network = "publicgoods"

	// MetisNetworkID is the ID for the Metis network.
	MetisNetworkID NetworkID = 1088
	MetisNetwork Network = "metis"

	// PolygonzkEVMNetworkID is the ID for the Polygon zkEVM network.
	PolygonzkEVMNetworkID NetworkID = 1101
	PolygonzkEVMNetwork Network = "polygonzkevm"

	// MoonbeamNetworkID is the ID for the Moonbeam network.
	MoonbeamNetworkID NetworkID = 1284
	MoonbeamNetwork Network = "moonbeam"

	// C1MilkomedaNetworkID is the ID for the C1 Milkomeda network.
	C1MilkomedaNetworkID NetworkID = 2001
	C1MilkomedaNetwork Network = "c1milkomeda"

	// MantleNetworkID is the ID for the Mantle network.
	MantleNetworkID NetworkID = 5000
	MantleNetwork Network = "mantle"

	// ZetaNetworkID is the ID for the Zeta network.
	ZetaNetworkID NetworkID = 7000
	ZetaNetwork Network = "zeta"

	// CyberNetworkID is the ID for the Cyber network.
	CyberNetworkID NetworkID = 7560
	CyberNetwork Network = "cyber"

	// BaseNetworkID is the ID for the Base network.
	BaseNetworkID NetworkID = 8453
	BaseNetwork Network = "base"

	// GnosisChiadoNetworkID is the ID for the Gnosis Chiado network.
	GnosisChiadoNetworkID NetworkID = 10200
	GnosisChiadoNetwork Network = "gnosischiado"

	// HoleskyNetworkID is the ID for the Holesky network.
	HoleskyNetworkID NetworkID = 17000
	HoleskyNetwork Network = "holesky"

	// ArbitrumOneNetworkID is the ID for the Arbitrum One network.
	ArbitrumOneNetworkID NetworkID = 42161
	ArbitrumOneNetwork Network = "arbitrumone"

	// ArbitrumNovaNetworkID is the ID for the Arbitrum Nova network.
	ArbitrumNovaNetworkID NetworkID = 42170
	ArbitrumNovaNetwork Network = "arbitrumnova"

	// CeloNetworkID is the ID for the Celo network.
	CeloNetworkID NetworkID = 42220
	CeloNetwork Network = "celo"

	// AvalancheNetworkID is the ID for the Avalanche network.
	AvalancheNetworkID NetworkID = 43114
	AvalancheNetwork Network = "avalanche"

	// LineaNetworkID is the ID for the Linea network.
	LineaNetworkID NetworkID = 59144
	LineaNetwork Network = "linea"

	// AmoyNetworkID is the ID for the Amoy network.
	AmoyNetworkID NetworkID = 80002
	AmoyNetwork Network = "amoy"

	// BlastNetworkID is the ID for the Blast network.
	BlastNetworkID NetworkID = 81457
	BlastNetwork Network = "blast"

	// BaseSepoliaNetworkID is the ID for the Base Sepolia network.
	BaseSepoliaNetworkID NetworkID = 84532
	BaseSepoliaNetwork Network = "basesepolia"

	// TaikoJolnrNetworkID is the ID for the Taiko Jolnr network.
	TaikoJolnrNetworkID NetworkID = 167007
	TaikoJolnrNetwork Network = "taikojolnr"

	// ArbitrumSepoliaNetworkID is the ID for the Arbitrum Sepolia network.
	ArbitrumSepoliaNetworkID NetworkID = 421614
	ArbitrumSepoliaNetwork Network = "arbitrumsepolia"

	// ScrollNetworkID is the ID for the Scroll network.
	ScrollNetworkID NetworkID = 534352
	ScrollNetwork Network = "scroll"

	// ZoraNetworkID is the ID for the Zora network.
	ZoraNetworkID NetworkID = 7777777
	ZoraNetwork Network = "zora"

	// SepoliaNetworkID is the ID for the Sepolia network.
	SepoliaNetworkID NetworkID = 11155111
	SepoliaNetwork Network = "sepolia"

	// OptimismSepoliaNetworkID is the ID for the Optimism Sepolia network.
	OptimismSepoliaNetworkID NetworkID = 11155420
	OptimismSepoliaNetwork Network = "optimismsepolia"

	// BlastSepoliaNetworkID is the ID for the Blast Sepolia network.
	BlastSepoliaNetworkID NetworkID = 168587773
	BlastSepoliaNetwork Network = "blastsepolia"

	// NeonEVMNetworkID is the ID for the Neon EVM network.
	NeonEVMNetworkID NetworkID = 245022934
	NeonEVMNetwork Network = "neonevm"

	// AuroraNetworkID is the ID for the Aurora network.
	AuroraNetworkID NetworkID = 1313161554
	AuroraNetwork Network = "aurora"

	// HarmonyShard0NetworkID is the ID for the Harmony Shard 0 network.
	HarmonyShard0NetworkID NetworkID = 1666600000
	HarmonyShard0Network Network = "harmonyshard0"

	// HarmonyShard1NetworkID is the ID for the Harmony Shard 1 network.
	HarmonyShard1NetworkID NetworkID = 1666600001
	HarmonyShard1Network Network = "harmonyshard1"
)

// String returns the string representation of the NetworkID.
func (n NetworkID) String() string {
	return n.ToBig().String()
}

// IsValid checks if the NetworkID is valid (non-zero).
func (n NetworkID) IsValid() bool {
	return n != 0
}


// ToNetwork converts a NetworkID to its corresponding Network.
func (n NetworkID) ToNetwork() Network {
	switch n {
	case EthereumNetworkID:
		return EthereumNetwork
	case GoerliNetworkID:
		return GoerliNetwork
	case OptimismNetworkID:
		return OptimismNetwork
	case FlareNetworkID:
		return FlareNetwork
	case RootstockNetworkID:
		return RootstockNetwork
	case LuksoNetworkID:
		return LuksoNetwork
	case CrabNetworkID:
		return CrabNetwork
	case DarwiniaNetworkID:
		return DarwiniaNetwork
	case BscNetworkID:
		return BscNetwork
	case GnosisNetworkID:
		return GnosisNetwork
	case PolygonNetworkID:
		return PolygonNetwork
	case ShimmerEVMNetworkID:
		return ShimmerEVMNetwork
	case MantaNetworkID:
		return MantaNetwork
	case XLayerTestnetNetworkID:
		return XLayerTestnetNetwork
	case XLayerNetworkID:
		return XLayerNetwork
	case FantomNetworkID:
		return FantomNetwork
	case KromaNetworkID:
		return KromaNetwork
	case BobaNetworkID:
		return BobaNetwork
	case ZksyncEraNetworkID:
		return ZksyncEraNetwork
	case PublicGoodsNetworkID:
		return PublicGoodsNetwork
	case MetisNetworkID:
		return MetisNetwork
	case PolygonzkEVMNetworkID:
		return PolygonzkEVMNetwork
	case MoonbeamNetworkID:
		return MoonbeamNetwork
	case C1MilkomedaNetworkID:
		return C1MilkomedaNetwork
	case MantleNetworkID:
		return MantleNetwork
	case ZetaNetworkID:
		return ZetaNetwork
	case CyberNetworkID:
		return CyberNetwork
	case BaseNetworkID:
		return BaseNetwork
	case GnosisChiadoNetworkID:
		return GnosisChiadoNetwork
	case HoleskyNetworkID:
		return HoleskyNetwork
	case ArbitrumOneNetworkID:
		return ArbitrumOneNetwork
	case ArbitrumNovaNetworkID:
		return ArbitrumNovaNetwork
	case CeloNetworkID:
		return CeloNetwork
	case AvalancheNetworkID:
		return AvalancheNetwork
	case LineaNetworkID:
		return LineaNetwork
	case AmoyNetworkID:
		return AmoyNetwork
	case BlastNetworkID:
		return BlastNetwork
	case BaseSepoliaNetworkID:
		return BaseSepoliaNetwork
	case TaikoJolnrNetworkID:
		return TaikoJolnrNetwork
	case ArbitrumSepoliaNetworkID:
		return ArbitrumSepoliaNetwork
	case ScrollNetworkID:
		return ScrollNetwork
	case ZoraNetworkID:
		return ZoraNetwork
	case SepoliaNetworkID:
		return SepoliaNetwork
	case OptimismSepoliaNetworkID:
		return OptimismSepoliaNetwork
	case BlastSepoliaNetworkID:
		return BlastSepoliaNetwork
	case NeonEVMNetworkID:
		return NeonEVMNetwork
	case AuroraNetworkID:
		return AuroraNetwork
	case HarmonyShard0NetworkID:
		return HarmonyShard0Network
	case HarmonyShard1NetworkID:
		return HarmonyShard1Network
	default:
		return EthereumNetwork
	}
}

// Uint64 converts a NetworkID to a uint64.
func (n NetworkID) Uint64() uint64 {
	return uint64(n)
}

// ToBig converts a NetworkID to a big.Int.
func (n NetworkID) ToBig() *big.Int {
	return new(big.Int).SetUint64(uint64(n))
}
