package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/osmosis-labs/osmosis/x/tokenfactory/keeper"
	"github.com/osmosis-labs/osmosis/x/tokenfactory/types"
)

func (suite *KeeperTestSuite) TestAdminMsgs() {
	suite.SetupTest()

	addr1 := sdk.AccAddress([]byte("addr1---------------"))
	addr2 := sdk.AccAddress([]byte("addr2---------------"))

	msgServer := keeper.NewMsgServerImpl(suite.app.TokenFactoryKeeper)

	// Create a denom
	res, err := msgServer.CreateDenom(sdk.WrapSDKContext(suite.ctx), types.NewMsgCreateDenom(addr1.String(), "bitcoin"))
	suite.Require().NoError(err)
	denom := res.GetNewTokenDenom()

	// Make sure that the admin is set correctly
	queryRes, err := suite.queryClient.DenomAuthorityMetadata(suite.ctx.Context(), &types.QueryDenomAuthorityMetadataRequest{
		Denom: res.GetNewTokenDenom(),
	})
	suite.Require().NoError(err)
	suite.Require().Equal(addr1.String(), queryRes.AuthorityMetadata.Admin)

	// Test minting to admins own account
	_, err = msgServer.Mint(sdk.WrapSDKContext(suite.ctx), types.NewMsgMint(addr1.String(), sdk.NewInt64Coin(denom, 10), addr1.String()))
	suite.Require().NoError(err)
	suite.Require().True(suite.app.BankKeeper.GetBalance(suite.ctx, addr1, denom).IsEqual(sdk.NewInt64Coin(denom, 10)), suite.app.BankKeeper.GetBalance(suite.ctx, addr1, denom))

	// Test minting to a separate account
	_, err = msgServer.Mint(sdk.WrapSDKContext(suite.ctx), types.NewMsgMint(addr1.String(), sdk.NewInt64Coin(denom, 10), addr2.String()))
	suite.Require().NoError(err)
	suite.Require().True(suite.app.BankKeeper.GetBalance(suite.ctx, addr2, denom).IsEqual(sdk.NewInt64Coin(denom, 10)))

	// Test force transferring
	_, err = msgServer.ForceTransfer(sdk.WrapSDKContext(suite.ctx), types.NewMsgForceTransfer(addr1.String(), sdk.NewInt64Coin(denom, 5), addr2.String(), addr1.String()))
	suite.Require().NoError(err)
	suite.Require().True(suite.app.BankKeeper.GetBalance(suite.ctx, addr1, denom).IsEqual(sdk.NewInt64Coin(denom, 15)))
	suite.Require().True(suite.app.BankKeeper.GetBalance(suite.ctx, addr2, denom).IsEqual(sdk.NewInt64Coin(denom, 5)))

	// Test burning
	_, err = msgServer.Burn(sdk.WrapSDKContext(suite.ctx), types.NewMsgBurn(addr1.String(), sdk.NewInt64Coin(denom, 5), addr2.String()))
	suite.Require().NoError(err)
	suite.Require().True(suite.app.BankKeeper.GetBalance(suite.ctx, addr2, denom).IsEqual(sdk.NewInt64Coin(denom, 0)))

	// Test Change Admin
	_, err = msgServer.ChangeAdmin(sdk.WrapSDKContext(suite.ctx), types.NewMsgChangeAdmin(addr1.String(), denom, addr2.String()))
	queryRes, err = suite.queryClient.DenomAuthorityMetadata(suite.ctx.Context(), &types.QueryDenomAuthorityMetadataRequest{
		Denom: res.GetNewTokenDenom(),
	})
	suite.Require().NoError(err)
	suite.Require().Equal(addr2.String(), queryRes.AuthorityMetadata.Admin)

	// Make sure can no longer do actions
	_, err = msgServer.Burn(sdk.WrapSDKContext(suite.ctx), types.NewMsgBurn(addr1.String(), sdk.NewInt64Coin(denom, 15), addr1.String()))
	suite.Require().Error(err)

	// Make sure the new admin works
	_, err = msgServer.Burn(sdk.WrapSDKContext(suite.ctx), types.NewMsgBurn(addr2.String(), sdk.NewInt64Coin(denom, 15), addr1.String()))
	suite.Require().NoError(err)
	suite.Require().True(suite.app.BankKeeper.GetBalance(suite.ctx, addr1, denom).IsEqual(sdk.NewInt64Coin(denom, 0)))
}
