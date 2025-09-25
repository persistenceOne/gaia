package keeper

import (
	"context"

	storetypes "cosmossdk.io/core/store"
	"cosmossdk.io/log"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/gaia/v24/x/liquid/types"
)

// Keeper of the x/liquid store
type Keeper struct {
	storeService  storetypes.KVStoreService
	tStoreService storetypes.TransientStoreService
	cdc           codec.BinaryCodec
	authKeeper    types.AccountKeeper
	bankKeeper    types.BankKeeper
	stakingKeeper types.StakingKeeper
	distKeeper    types.DistributionKeeper
	authority     string
}

// NewKeeper creates a new liquid Keeper instance
func NewKeeper(
	cdc codec.BinaryCodec,
	storeService storetypes.KVStoreService,
	tStoreService storetypes.TransientStoreService,
	ak types.AccountKeeper,
	bk types.BankKeeper,
	sk types.StakingKeeper,
	dk types.DistributionKeeper,
	authority string,
) *Keeper {
	// ensure that authority is a valid AccAddress
	if _, err := ak.AddressCodec().StringToBytes(authority); err != nil {
		panic("authority is not a valid acc address")
	}

	return &Keeper{
		storeService:  storeService,
		tStoreService: tStoreService,
		cdc:           cdc,
		authKeeper:    ak,
		bankKeeper:    bk,
		stakingKeeper: sk,
		distKeeper:    dk,
		authority:     authority,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx context.Context) log.Logger {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	return sdkCtx.Logger().With("module", "x/"+types.ModuleName)
}

// GetAuthority returns the x/liquid module's authority.
func (k Keeper) GetAuthority() string {
	return k.authority
}

// GetTransientHookDelegationData retrieves delegation data from the transient store using a predefined key.
// Returns the delegation and an error if unmarshalling or store retrieval fails.
func (k Keeper) GetTransientHookDelegationData(ctx context.Context) (*stakingtypes.Delegation, error) {
	tStore := k.tStoreService.OpenTransientStore(ctx)
	delBz, err := tStore.Get(types.TransientHookDelegationDataPrefix)
	if err != nil || delBz == nil {
		return nil, err
	}
	var del stakingtypes.Delegation
	err = k.cdc.Unmarshal(delBz, &del)
	if err != nil {
		return nil, err
	}
	return &del, nil
}

// SetTransientHookDelegationData sets delegation data into the transient store using serialized input for temporary storage.
func (k Keeper) SetTransientHookDelegationData(ctx context.Context, del *stakingtypes.Delegation) error {
	tStore := k.tStoreService.OpenTransientStore(ctx)
	delbBz, err := k.cdc.Marshal(del)
	if err != nil {
		return err
	}
	return tStore.Set(types.TransientHookDelegationDataPrefix, delbBz)
}

// ResetTransientHookDelegationData clears transient hook delegation data from the transient store for the given context.
func (k Keeper) ResetTransientHookDelegationData(ctx context.Context) error {
	tStore := k.tStoreService.OpenTransientStore(ctx)
	return tStore.Delete(types.TransientHookDelegationDataPrefix)

}
