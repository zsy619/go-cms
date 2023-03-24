layui.define(['form', 'jquery'], function (exports) {
    var $ = layui.jquery;

    layui.link(layui.cache.base + 'treeSelect2/treeSelect2.css');

    var treeSelected = 'layui-form-selected',
        nodeSelected = 'curSelectedNode',
        TREE_INPUT_CLASS = 'layui-treeselect',
        TREE_SELECT_CLASS = 'layui-treeSelect',
        TREE_SELECT_TITLE_CLASS = 'layui-select-title',
        TREE_SELECT_BODY_CLASS = 'layui-treeSelect-body',
        TREE_SELECT_SEARCHED_CLASS = 'layui-treeSelect-search-ed';

    var obj = {
        get: function (filter) {
            return $('*[lay-filter=' + filter + ']');
        },
        filter: function (filter) {
            var tf = this.get(filter),
                o = tf.next();
            return o;
        },
        treeObj: function (filter) {
            var o = this.filter(filter),
                treeId = o.find('.' + TREE_SELECT_BODY_CLASS).attr('id'),
                tree = $.fn.zTree.getZTreeObj(treeId);
            return tree;
        }
    };

    var treeSelect2 = {
        render: function (options) {
            var opt = $.extend({
                type: "get",
                search: false,
                multipleSelect: false,
                headers: [],
                valueId: 'id'
            }, options);
            var t = new Date().getTime(),
                TREE_INPUT_ID = 'treeSelect-input-' + t,
                TREE_SELECT_ID = 'layui-treeSelect-' + t,
                TREE_SELECT_TITLE_ID = 'layui-select-title-' + t,
                TREE_SELECT_BODY_ID = 'layui-treeSelect-body-' + t;
            var selector = {
                TREE_INPUT_ID: TREE_INPUT_ID,
                TREE_SELECT_ID: TREE_SELECT_ID,
                TREE_SELECT_TITLE_ID: TREE_SELECT_TITLE_ID,
                TREE_SELECT_BODY_ID: TREE_SELECT_BODY_ID,
                options: opt,
                onCollapse: function () {
                    selector.focusInput();
                },
                onExpand: function () {
                    selector.focusInput();
                },
                beforeExpand: function () {
                },
                focusInput: function () {
                    $('#' + TREE_INPUT_ID).focus();
                },
                onClick: function (event, treeId, treeNode) {
                    var tree = $.fn.zTree.getZTreeObj(treeId);
                    var name = treeNode[tree.setting.data.key.name],
                        id = treeNode[tree.selector.options.valueId],
                        $input = $('#' + TREE_SELECT_TITLE_ID + ' input');
                    $input.val(name);
                    $(opt.elem).attr('value', id).val(id);
                    $('#' + TREE_SELECT_ID).removeClass(treeSelected);

                    if (opt.click) {
                        var obj = {
                            data: selector.data,
                            current: treeNode,
                            treeId: TREE_SELECT_ID
                        };
                        opt.click(obj);
                    }
                    return this;
                },
                setting: function () {
                    opt.treeSetting = $.extend(opt.treeSetting, {
                        callback: {
                            onClick: this.onClick,
                            onExpand: this.onExpand,
                            onCollapse: this.onCollapse,
                            beforeExpand: this.beforeExpand
                        }
                    });
                    return opt.treeSetting;
                },
                init: function () {
                    this.hideElem().create().toggleSelect().preventEvent();
                    $.fn.zTree.init($('#' + TREE_SELECT_BODY_ID), this.setting(), this.options.data);
                    var treeObj = $.fn.zTree.getZTreeObj(TREE_SELECT_BODY_ID);
                    treeObj.selector = this;
                    if (opt.success) {
                        var obj = {
                            treeId: TREE_SELECT_ID,
                            data: this.options.data
                        };
                        opt.success(obj);
                    }
                },
                refresh: function (options) {
                    var treeObj = $.fn.zTree.getZTreeObj(this.TREE_SELECT_BODY_ID);
                    this.options = $.extend(this.options, options);
                    treeObj.destroy();
                    $.fn.zTree.init($('#' + this.TREE_SELECT_BODY_ID), this.options.treeSetting, this.options.data);
                    var treeObj = $.fn.zTree.getZTreeObj(this.TREE_SELECT_BODY_ID);
                    treeObj.selector = this;
                    if (opt.success) {
                        var obj = {
                            treeId: TREE_SELECT_ID,
                            data: this.options.data
                        };
                        opt.success(obj);
                    }
                },
                hideElem: function () {
                    $(opt.elem).hide();
                    return this;
                },
                create: function () {
                    var readonly = '';
                    if (!opt.search) {
                        readonly = 'readonly';
                    }
                    var placeholder = $(opt.elem).attr('placeholder') || "";
                    var selectHtml = '<div class="' + TREE_SELECT_CLASS + ' layui-unselect layui-form-select" id="' + TREE_SELECT_ID + '">' +
                        '<div class="' + TREE_SELECT_TITLE_CLASS + '" id="' + TREE_SELECT_TITLE_ID + '">' +
                        ' <input type="text" id="' + TREE_INPUT_ID + '" placeholder="' + placeholder + '" value="" ' + readonly + ' class="layui-input layui-unselect">' +
                        '<i class="layui-edge"></i>' +
                        '</div>' +
                        '<div class="layui-anim layui-anim-upbit" style="">' +
                        '<div class="' + TREE_SELECT_BODY_CLASS + ' ztree" id="' + TREE_SELECT_BODY_ID + '"></div>' +
                        '</div>' +
                        '</div>';
                    $(opt.elem).parent().append(selectHtml);
                    return this;
                },
                toggleSelect: function () {
                    var item = '#' + TREE_SELECT_TITLE_ID;
                    selector.event('click', item, function (e) {
                        var $select = $('#' + TREE_SELECT_ID);
                        if ($select.hasClass(treeSelected)) {
                            $select.removeClass(treeSelected);
                            $('#' + TREE_INPUT_ID).blur();
                        } else {
                            // 隐藏其他picker
                            $('.layui-form-select').removeClass(treeSelected);
                            // 显示当前picker
                            $select.addClass(treeSelected);
                        }
                        e.stopPropagation();
                    });
                    $(document).click(function () {
                        var $select = $('#' + TREE_SELECT_ID);
                        if ($select.hasClass(treeSelected)) {
                            $select.removeClass(treeSelected);
                            $('#' + TREE_INPUT_ID).blur();
                        }
                    });
                    return this;
                },
                preventEvent: function () {
                    var item = '#' + TREE_SELECT_ID + ' .layui-anim';
                    this.event('click', item, function (e) {
                        e.stopPropagation();
                    });
                    return this;
                },
                event: function (evt, el, fn) {
                    $('body').on(evt, el, fn);
                }
            };
            //如果已经构建过
            var el = $(opt.elem);
            var treeObj = obj.treeObj(el.attr('lay-filter'));

            if (opt.data == null) {
                throw 'argument invalid "data"';
            } else if (typeof (opt.data) === "string") {
                $.ajax({
                    url: opt.data,
                    type: opt.type,
                    headers: opt.headers,
                    dataType: 'json',
                    success: function (d) {
                        if (treeObj != null) {
                            treeObj.selector.data = d;
                            treeObj.selector.refresh();
                        } else {
                            selector.data = d;
                            selector.init();
                        }
                    }
                });
            } else {
                if (treeObj != null) {
                    treeObj.selector.data = opt.data;
                    treeObj.selector.refresh();
                } else {
                    selector.data = opt.data;
                    selector.init();
                }
            }
            if (treeObj != null) {
                return treeObj.selector;
            } else {
                return selector;
            }
        },
        refresh: function (filter, options) {
            var o = obj.filter(filter);
            var treeObj = obj.treeObj(filter);
            var opt = treeObj.selector.options;
            if (options && options.data) {
                opt.data = options.data;
            }
            this.render(opt);
        },
        disableNode: function (filter, id, disabled) {
            var o = obj.filter(filter),
                treeInput = o.find('.' + TREE_SELECT_TITLE_CLASS + ' input'),
                treeObj = obj.treeObj(filter),
                node = treeObj.getNodeByParam(treeObj.selector.options.valueId, id, null),
                name = node[treeObj.setting.data.key.name];
            treeInput.val(name);
            o.find('a[treenode_a]').removeClass(nodeSelected);
            obj.get(filter).val(id).attr('value', id);
            if (treeObj.setting.check.enable) {
                treeObj.setChkDisabled(node, disabled);
            } else {
                if (disabled) {
                    treeObj.hideNode(node);
                } else {
                    treeObj.showNode(node);
                }
            }
        },
        checkNode: function (filter, ids, checked) {
            var o = obj.filter(filter),
                treeInput = o.find('.' + TREE_SELECT_TITLE_CLASS + ' input'),
                treeObj = obj.treeObj(filter);
            if (treeObj.setting.check.enable) {
                for (var i = 0; i < ids.length; i++) {
                    var id = ids[i];
                    var node = treeObj.getNodeByParam(treeObj.selector.options.valueId, id, null),
                        name = node[treeObj.setting.data.key.name];
                    treeInput.val(name);
                    o.find('a[treenode_a]').removeClass(nodeSelected);
                    obj.get(filter).val(id).attr('value', id);
                    treeObj.checkNode(node, checked, true);
                }
            }
        },
        selectNode: function (filter, id) {
            var o = obj.filter(filter),
                treeInput = o.find('.' + TREE_SELECT_TITLE_CLASS + ' input'),
                treeObj = obj.treeObj(filter),
                node = treeObj.getNodeByParam(treeObj.selector.options.valueId, id, null),
                name = node[treeObj.setting.data.key.name];
            treeInput.val(name);
            o.find('a[treenode_a]').removeClass(nodeSelected);
            obj.get(filter).val(id).attr('value', id);
            treeObj.selectNode(node);
        },
        revokeNode: function (filter, fn) {
            var o = obj.filter(filter);
            o.find('a[treenode_a]').removeClass(nodeSelected);
            o.find('.' + TREE_SELECT_TITLE_CLASS + ' input.layui-input').val('');
            obj.get(filter).attr('value', '').val('');
            obj.treeObj(filter).expandAll(false);
            if (fn) {
                fn({
                    treeId: o.attr('id')
                });
            }
        },
        destroy: function (filter) {
            var o = obj.filter(filter);
            o.remove();
            var treeObj = obj.treeObj(filter);
            treeObj.destroy();
            obj.get(filter).show();
        },
        zTree: function (filter) {
            return obj.treeObj(filter);
        }
    };

    //输出接口
    exports('treeSelect2', treeSelect2);
});    