const copyToClipboard = (target = "") => {
    if (!target) throw new Error("没有输出可复制");

    navigator.clipboard
        ?.writeText(target)
        .then(() => {})
        .catch(() => {
            // 回退
            const ta = document.createElement("textarea");
            ta.value = target;
            ta.style.position = "fixed";
            ta.style.opacity = "0";
            document.body.appendChild(ta);
            ta.select();
            try {
                document.execCommand("copy");
            } catch {
                throw new Error("复制失败");
            } finally {
                document.body.removeChild(ta);
            }
        });
};

export default { copyToClipboard };
