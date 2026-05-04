$('.theme-box .theme-item').click(function(e) {
	var skin = e.target.dataset.theme;
	localStorage.setItem('skin', skin);
	$('.skin').remove();
  	location.reload();
	renderTheme(skin)
})

var curSkin = localStorage.getItem('skin')
$('.'+curSkin).addClass('active')
//init(curSkin)
function init(skin){
	renderPagination('skin_1','#0B7C81')  
}
renderTheme(curSkin)
function renderTheme(skin){
	var skin = localStorage.getItem('skin');
	console.log(skin)
	//skin_0:blue,
	//skin_1:green,
	//skin_2:red	
	if(skin == 'skin_0' || skin==undefined){
		$('.skin_0').addClass('active')
		$('head').append('<link href="/views/themes/subIV/static/css/skin_0.css?v='+Math.random()+'" rel="stylesheet" type="text/css" class="skin"/>')
		renderPagination('skin_0','#285396')  
	}
	if(skin == 'skin_1'){
		$('head').append('<link href="/views/themes/subIV/static/css/skin_1.css?v='+Math.random()+'" rel="stylesheet" type="text/css" class="skin"/>')
		renderPagination('skin_1','#0B7C81')  
	}
	if(skin == 'skin_2'){
		$('head').append('<link href="/views/themes/subIV/static/css/skin_2.css?v='+Math.random()+'" rel="stylesheet" type="text/css" class="skin"/>/>')
		renderPagination('skin_2','#5A0000')  
	}
}



function renderPagination(skin,color){
	layui.use(['laypage'], function() {
		var laypage = layui.laypage;
		laypage.render({
			elem: 'pagination',
			count: 100,
			theme: color
		});
	})
}
function getUrlParams() {
	var url = location.search;
	var params = new Object(); 
	if(url.indexOf("?") != -1) { //判断 URL 中是否包含查询字符串
		var str = url.substr(1); //如果 URL 中包含查询字符串，截取查询字符串，去掉前面的“?”号。
		strs = str.split("&"); //将查询字符串按“&”号分割成一个个参数对。
		for(var i = 0; i < strs.length; i++) { //循环遍历所有的参数对。
			params[strs[i].split("=")[0]] = unescape(strs[i].split("=")[1]);
		}
	}
	return params;
}

toggleNav()
function toggleNav(){
	$('.app-nav-btn').click(function(){
		$('.app-nav-box').toggleClass('show')
		$('.nav-mask').show()
		$('body').addClass('fixed')
	})
	$('.app-nav-box .close').click(function(){
		$('.app-nav-box').removeClass('show')
		$('.nav-mask').hide()
		$('body').removeClass('fixed')
	})
	$('.app-nav-box i').click(function(){
		var li = $(this).parent('li')
		
		if(li.hasClass('on')){
			li.children('.list').slideUp(600);
			li.removeClass('on')
		}else{
			li.children('.list').slideDown();
			li.addClass('on')
		}
	})
}
