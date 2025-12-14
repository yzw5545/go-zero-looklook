package user

import (
	"context"

	"looklook/app/usercenter/cmd/api/internal/svc"
	"looklook/app/usercenter/cmd/api/internal/types"
	"looklook/app/usercenter/cmd/rpc/usercenter"
	"looklook/pkg/ctxdata"
	"looklook/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type ChangePasswordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// change password
func NewChangePasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChangePasswordLogic {
	return &ChangePasswordLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ChangePasswordLogic) ChangePassword(req *types.ChangePasswordReq) (resp *types.ChangePasswordResp, err error) {
	// 1. 参数验证
	if req.OldPassword == "" {
		return nil, xerr.NewErrMsg("旧密码不能为空")
	}

	if req.NewPassword == "" {
		return nil, xerr.NewErrMsg("新密码不能为空")
	}

	if req.ConfirmPassword == "" {
		return nil, xerr.NewErrMsg("确认密码不能为空")
	}

	// 2. 验证新密码和确认密码是否一致
	if req.NewPassword != req.ConfirmPassword {
		return nil, xerr.NewErrMsg("两次输入的新密码不一致")
	}

	// 3. 验证新密码和旧密码不能相同
	if req.NewPassword == req.OldPassword {
		return nil, xerr.NewErrMsg("新密码不能与旧密码相同")
	}

	// 4. 验证密码长度
	if len(req.NewPassword) < 6 || len(req.NewPassword) > 20 {
		return nil, xerr.NewErrMsg("密码长度必须在6-20位之间")
	}

	// 验证密码复杂度（必须包含字母和数字）
	hasLetter := false
	hasNumber := false
	for _, c := range req.NewPassword {
		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') {
			hasLetter = true
		} else if c >= '0' && c <= '9' {
			hasNumber = true
		}
	}
	if !hasLetter || !hasNumber {
		return nil, xerr.NewErrMsg("密码必须包含字母和数字")
	}

	// 5. 获取当前用户ID（从JWT token中）
	userId := ctxdata.GetUidFromCtx(l.ctx)

	// 6. 调用RPC服务修改密码
	_, err = l.svcCtx.UsercenterRpc.ChangePassword(l.ctx, &usercenter.ChangePasswordReq{
		UserId:      userId,
		OldPassword: req.OldPassword,
		NewPassword: req.NewPassword,
	})

	if err != nil {
		return nil, err
	}

	return &types.ChangePasswordResp{}, nil
}

