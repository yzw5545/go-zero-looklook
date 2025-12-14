package logic

import (
	"context"

	"looklook/app/usercenter/cmd/rpc/internal/svc"
	"looklook/app/usercenter/cmd/rpc/pb"
	"looklook/app/usercenter/model"
	"looklook/pkg/tool"
	"looklook/pkg/xerr"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type ChangePasswordLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

var ErrOldPasswordError = xerr.NewErrMsg("旧密码错误")

func NewChangePasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChangePasswordLogic {
	return &ChangePasswordLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ChangePasswordLogic) ChangePassword(in *pb.ChangePasswordReq) (*pb.ChangePasswordResp, error) {
	// 1. 根据用户ID查询用户信息
	user, err := l.svcCtx.UserModel.FindOne(l.ctx, in.UserId)
	if err != nil && err != model.ErrNotFound {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_ERROR), "查询用户信息失败，userId:%d,err:%v", in.UserId, err)
	}
	if user == nil {
		return nil, errors.Wrapf(xerr.NewErrMsg("用户不存在"), "userId:%d", in.UserId)
	}

	// 2. 验证旧密码是否正确
	if tool.Md5ByString(in.OldPassword) != user.Password {
		return nil, errors.Wrapf(ErrOldPasswordError, "旧密码验证失败, userId:%d", in.UserId)
	}

	// 3. 加密新密码
	newPasswordMd5 := tool.Md5ByString(in.NewPassword)

	// 4. 更新数据库中的密码
	user.Password = newPasswordMd5
	_, err = l.svcCtx.UserModel.Update(l.ctx, nil, user)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_ERROR), "更新密码失败，userId:%d,err:%v", in.UserId, err)
	}

	return &pb.ChangePasswordResp{}, nil
}
