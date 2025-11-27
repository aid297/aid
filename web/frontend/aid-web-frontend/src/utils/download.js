const saveToFile = (target = "", filename = "", type = "text/plain;charset=utf-8") => {
    if (!target) {
        throw new Error("下载内容不能为空");
    }

    const blob = new Blob([target], { type });
    const url = URL.createObjectURL(blob);
    const tmp = document.createElement("a");
    tmp.href = url;
    tmp.download = filename || "download.txt";
    document.body.appendChild(tmp);
    tmp.click();
    document.body.removeChild(tmp);
    URL.revokeObjectURL(url);
};

export default { saveToFile };
