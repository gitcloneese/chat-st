package logic

import (
	pb "xy3-proto/logger"
)

func (p *Logic) initRoute() {
	//TODO: 示例
	p.addRoute(pb.ELogCategory_ELC_Login, "Login")
	p.addRoute(pb.ELogCategory_ELC_Register, "Register")
	p.addRoute(pb.ELogCategory_ELC_Logout, "Logout")
	p.addRoute(pb.ELogCategory_ELC_Resource, "Resource")
	p.addRoute(pb.ELogCategory_ELC_Mail, "Mail")
	p.addRoute(pb.ELogCategory_ELC_Online, "Online")
	p.addRoute(pb.ELogCategory_ELC_Task, "Task")
	p.addRoute(pb.ELogCategory_ELC_Battle, "Battle")
	p.addRoute(pb.ELogCategory_ELC_PVP_Ranking, "PvpRanking")
	p.addRoute(pb.ELogCategory_ELC_Tutorial_Progress, "TutorialProgress")
	p.addRoute(pb.ELogCategory_ELC_OnlineTime, "OnlineTime")
}
