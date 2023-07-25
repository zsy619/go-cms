document.write(`<div class="header">
			<div class="container clearfix">
				<div class="logo">
					<img src="../static/images/logo_img.png" alt="">
					<span class="title">郑州大学就业信息服务平台</span>
				</div>
				<!--登录前-->
				<!--<div class="user-box">
					<a href="stu-1.1-个人中心-个人信息.html" class="one">求职者登录/注册</a>
					<a href="javascript:void(0)" class="two">用人单位登录/注册</a>
				</div>-->
				<!--登录后-->
				<div class="theme-box">
					<div class="theme-item skin_0" data-theme="skin_0"></div>
					<div class="theme-item skin_1" data-theme="skin_1"></div>
					<div class="theme-item skin_2" data-theme="skin_2"></div>
				</div>
				<div class="user-box">
					<a href="javascript:void(0)" class="user">王博与
						<img src="../static/images/default_stu.png" alt="" />
					</a>
					<div class="sel-list">
						<div class="m-arr"></div>
						<a href="comp-6.3-个人中心-企业主页.html">企业主页</a>
						<a href="comp-6.2-个人中心-账号设置.html">账号</a>
						<a href="javascript:void(0)">退出</a>
					</div>
				</div>
				<div class="app-nav-btn">
					<img src="../static/images/icon_btn-nav.png"/>
				</div>
			</div>
		</div>
		<div class="nav-box">
			<div class="container">
				<ul class="clearfix head-nav">
					<li class="active">
						<a href="index.html">首&nbsp;&nbsp;&nbsp;&nbsp;页</a>
					</li>
					<li>
						<a href="#">求职者</a>
						<div class="subNav">
							<a href="javascript:void(0)">项目通知</a>
							<a href="javascript:void(0)">专项研究项目</a>
							<a href="javascript:void(0)">示范校建设项目</a>
							<a href="javascript:void(0)">研究实验室项目</a>
							<a href="javascript:void(0)">教改案例项目</a>
							<a href="javascript:void(0)">……</a>
						</div>
					</li>
					<li>
						<a href="#">用人单位</a>
					</li>
					<li>
						<a href="">创业天地 </a>
						<div class="subNav">
							<a href="javascript:void(0)">新 闻</a>
							<a href="javascript:void(0)">政 策</a>
							<a href="javascript:void(0)">行业动态</a>
						</div>
					</li>
					<li>
						<a href="下载中心.html">下载中心</a>
					</li>
					<li>
						<a href="">关于我们</a>
						<div class="subNav">
							<a href="javascript:void(0)">企业白名单</a>
							<a href="javascript:void(0)">产品白名单</a>
						</div>
					</li>
					
				</ul>
			</div>
		</div>
		<div class="app-nav-box">
			<div class="other">
				<a href="#">求职者</a>
				<span class="line">|</span>
				<a href="#">用人单位</a>
				<span class="line">|</span>
				<a href="#">学校</a>
			</div>
			<ul>
				<li>
					<a href="index.html">首页</a>
				</li>
				<li>
					<a href="#">学生</a>
					<i></i>
					<div class="list">
						<a href="#">项目通知</a>
						<a href="#">专项研究项目</a>
						<a href="#">示范校建设项目</a>
						<a href="#">研究实验室项目</a>
					</div>
				</li>
				<li>
					<a href="#">就业单位</a>
				</li>
				<li>
					<a href="">创业天地 </a>
					<i></i>
					<div class="list">
						<a href="javascript:void(0)">新 闻</a>
						<a href="javascript:void(0)">政 策</a>
						<a href="javascript:void(0)">行业动态</a>
					</div>
				</li>
				<li>
					<a href="下载中心.html">下载中心</a>
				</li>
				<li>
					<a href="">关于我们</a>					
				</li>
			</ul>
			<div class="close">×</div>
		</div>
		`)


//function navActive(item){
//	$('.nav-box ul li').each(function(i,ele){
//		if($(ele).children('.nav-th').find('a').html() == item){
//			$(ele).addClass('active')
//		}
//	})
//}