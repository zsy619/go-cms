layui.use(['form'], function () {
    var form = layui.form;
    var layer = layui.layer;
    form.render();
    if (top.location !== self.location) { top.location = self.location; }
    $('.bind-password').on('click', function () {
        if ($(this).hasClass('icon-5')) {
            $(this).removeClass('icon-5');
            $('input[name="password"]').attr('type', 'password');
        } else {
            $(this).addClass('icon-5');
            $('input[name="password"]').attr('type', 'text');
        }
    });
    $('.icon-nocheck').on('click', function () {
        if ($(this).hasClass('icon-check')) { $(this).removeClass('icon-check'); }
        else { $(this).addClass('icon-check'); }
    });
    form.on('submit(login)', function (data) {
        data = data.field;
        if (data.username === '') { layer.msg('用户名不能为空'); return false; }
        if (data.password === '') { layer.msg('密码不能为空'); return false; }
        if (data.captcha === '') { layer.msg('验证码不能为空'); return false; }
        var psdRole = checkPasswordRole(data.password);
        if (psdRole.length > 0) { layer.msg(psdRole); return false; }
        $.ajax({
            url: '/cms/admin/login/verify',
            type: 'post',
            data: { username: data.username, password: data.password, captcha: data.captcha },
            beforeSend: function () { this.layerIndex = layer.load(0, { shade: [0.5, '#393D49'] }); },
            success: function (res) {
                if (res.code === 1) { layer.msg(res.message || res.msg, { icon: 5 }); return; }
                if (res.code === 0) {
                    layer.msg(res.message || res.msg, { icon: 6, time: 1000 }, function () {
                        window.location.href = '/admin/index';
                    });
                } else {
                    layer.msg(res.message || res.msg, { icon: 5 });
                    $('input[name="username"]').val('');
                    $('input[name="password"]').val('');
                }
            },
            error: function (xhr, status, error) {
                console.error('AJAX错误:', status, error);
                layer.msg('请求失败: ' + error, { icon: 5 });
                return false;
            },
            complete: function () {
                $('#captcha').val('');
                $('#refreshCaptcha').click();
                layer.close(this.layerIndex);
            }
        });
        return false;
    });
    function checkPasswordRole(password) {
        if (password.length < 8 || password.length > 16) { return '密码长度应为8~16位'; }
        var regex = /^(?=.*?[A-Z])(?=.*?[a-z])(?=.*?[0-9])(?=.*?[^\w\s])/;
        if (!regex.test(password)) { return '密码校验未通过'; }
        return '';
    }
});
