layui.use(['form'], function () {
    var form = layui.form,
        layer = layui.layer;

    // 登录过期的时候，跳出ifram框架
    if (top.location != self.location) top.location = self.location;

    // 粒子线条背景
    $(document).ready(function () {
        $('.layui-container').particleground({
            dotColor: '#7ec7fd',
            lineColor: '#7ec7fd'
        });
    });

    // 进行登录操作
    form.on('submit(login)', function (data) {
        data = data.field;
        if (data.username == '') {
            layer.msg('用户名不能为空');
            return false;
        }
        if (data.password == '') {
            layer.msg('密码不能为空');
            return false;
        }
        if (data.captcha == '') {
            layer.msg('验证码不能为空');
            return false;
        }
        // var loading = layer.msg('处理中，请稍后...', { icon: 16, shade: 0.3, time: 0 });
        $.ajax({
            url: "/cms/admin/login/verify",
            type: 'post',
            data: { username: data.username, password: data.password, captcha: data.captcha },
            beforeSend: function () {
                this.layerIndex = layer.load(0, { shade: [0.5, '#393D49'] });
            },
            success: function (data) {
                console.log(data);
                switch (data.code) {
                    case 1:
                        layer.msg(data.msg, { icon: 5 });//失败的表情
                        return;
                    case 0:
                        layer.msg(data.msg, {
                            icon: 6,//成功的表情
                            time: 1000 //1秒关闭（如果不配置，默认是3秒）
                        }, function () {
                            window.location = '/admin/index';
                        });
                        break;
                    default:
                        layer.msg(data.msg, { icon: 5 });//失败的表情
                        $("input[name=username]").val("");
                        $("input[name='password']").val("");
                        return;
                }
            },
            complete: function () {
                $("#captcha").val("");
                $("#refreshCaptcha").click();
                layer.close(this.layerIndex);
            },
        });
        return false;
    });
});