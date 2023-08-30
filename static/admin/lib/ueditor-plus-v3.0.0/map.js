// UE.registerUI('map', function (editor, uiName) {
//     //创建dialog
//     var dialog = new UE.ui.Dialog({
//         //指定弹出层中页面的路径，这里只能支持页面,因为跟addCustomizeDialog.js相同目录，所以无需加路径
//         iframeUrl: '/static/admin/lib/ueditor-plus-v3.0.0/dialogs/map/map.html',
//         //需要指定当前的编辑器实例
//         editor: editor,
//         //指定dialog的名字
//         name: uiName,
//         //dialog的标题
//         title: "百度地图",

//         //指定dialog的外围样式
//         cssRules: "width:600px;height:400px;",

//         //如果给出了buttons就代表dialog有确定和取消
//         buttons: [
//             {
//                 className: 'edui-okbutton',
//                 label: '确定',
//                 onclick: function () {
//                     dialog.close(true);
//                 }
//             },
//             {
//                 className: 'edui-cancelbutton',
//                 label: '取消',
//                 onclick: function () {
//                     dialog.close(false);
//                 }
//             }
//         ]
//     });

//     //参考addCustomizeButton.js
//     var btn = new UE.ui.Button({
//         name: 'dialogbutton' + uiName,
//         title: 'dialogbutton' + uiName,
//         //需要添加的额外样式，指定icon图标，这里默认使用一个重复的icon
//         cssRules: 'background-position: -500px 0;',
//         onclick: function () {
//             //渲染dialog
//             dialog.render();
//             dialog.open();
//         }
//     });

//     return btn;
// }/*index 指定添加到工具栏上的那个位置，默认时追加到最后,editorId 指定这个UI是那个编辑器实例上的，默认是页面上所有的编辑器都会添加这个按钮*/);


UE.plugins["map"] = function () {
	UE.commands["map"] = {
		execCommand: function (cmd, obj) {
			obj.html && this.execCommand("inserthtml", obj.html);
		}
	};
	this.addListener("click", function (type, evt) {
		var el = evt.target || evt.srcElement,
			range = this.selection.getRange();
		var tnode = UE.dom.domUtils.findParent(
			el,
			function (node) {
				if (node.className && UE.dom.domUtils.hasClass(node, "ue_t")) {
					return node;
				}
			},
			true
		);
		tnode && range.selectNode(tnode).shrinkBoundary().select();
	});
	this.addListener("keydown", function (type, evt) {
		var range = this.selection.getRange();
		if (!range.collapsed) {
			if (!evt.ctrlKey && !evt.metaKey && !evt.shiftKey && !evt.altKey) {
				var tnode = UE.dom.domUtils.findParent(
					range.startContainer,
					function (node) {
						if (node.className && UE.dom.domUtils.hasClass(node, "ue_t")) {
							return node;
						}
					},
					true
				);
				if (tnode) {
					UE.dom.domUtils.removeClasses(tnode, ["ue_t"]);
				}
			}
		}
	});
};
