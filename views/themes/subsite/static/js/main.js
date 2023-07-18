renderPagination(green,'#0B7C81')
$('.theme-box .theme-item').click(function(e) {
	
	var theme = e.target.dataset.theme;
	console.log(theme)
	$(this).addClass('active').siblings().removeClass('active')
	if(theme == 'blue') {
		document.documentElement.setAttribute('theme', '')
		renderPagination('#285396')
	}
	if(theme == 'green') {
		document.documentElement.setAttribute('theme', 'green')
		renderPagination('#0B7C81')
	}
	if(theme == 'red') {
		document.documentElement.setAttribute('theme', 'red')
		renderPagination('#5A0000')
	}

})

function renderPagination(color){
	layui.use(['laypage'], function() {
		var laypage = layui.laypage;
		laypage.render({
			elem: 'pagination',
			count: 100,
			theme: color
		});
	})
}
