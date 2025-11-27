<template>
    <q-page class="q-pa-md">
        <q-card flat square bordered class="my-card">
            <q-card-section>
                <div class="text-h4 text-deep-orange">Base64 编码</div>
            </q-card-section>

            <q-card-section class="q-pt-none">
                <div class="row q-gutter-md">
                    <div class="col">
                        <q-btn square label="文本 → Base64" color="primary" @click="encodeText" class="full-width" />
                    </div>
                    <div class="col">
                        <q-btn square label="Base64 → 文本" color="secondary" @click="decodeText" class="full-width" />
                    </div>
                    <div class="col">
                        <q-btn square label="复制输出" color="purple" @click="copyBase64Output" :disable="!base64OutputText" class="full-width" />
                    </div>
                    <div class="col">
                        <q-btn square label="下载为文件" color="amber" @click="downloadBase64Output" :disable="!base64OutputText" class="full-width" />
                    </div>
                    <div class="col">
                        <q-btn square label="清除" color="negative" @click="clearAll" class="full-width" />
                    </div>
                </div>
            </q-card-section>

            <q-separator inset />

            <q-card-section>
                <div class="row q-gutter-md q-mt-sm">
                    <div class="col">
                        <q-input v-model="base64InputText" type="textarea" autogrow outlined square label="输入文本 / Base64" :placeholder="placeholder" ref="txtBase64InputText" />
                    </div>
                    <div class="col">
                        <q-input v-model="base64OutputText" type="textarea" autogrow outlined square label="输出" readonly />
                        <div class="q-mt-sm text-caption">支持 UTF-8 编码。解码时会自动去除空白字符。</div>
                    </div>
                </div>
            </q-card-section>
        </q-card>

        <q-card flat square bordered class="my-card q-mt-md">
            <q-card-section>
                <div class="text-h4 text-deep-orange">哈希编码</div>
            </q-card-section>

            <q-card-section class="q-pt-none">
                <div class="row q-gutter-md">
                    <div class="col">
                        <q-btn square label="SHA1" color="primary" @click="encodeSHA1" class="full-width" />
                    </div>
                    <div class="col">
                        <q-btn square label="SHA224" color="primary" @click="encodeSHA224" class="full-width" />
                    </div>
                    <div class="col">
                        <q-btn square label="SHA256" color="primary" @click="encodeSHA256" class="full-width" />
                    </div>
                    <div class="col">
                        <q-btn square label="SHA384" color="primary" @click="encodeSHA384" class="full-width" />
                    </div>
                    <div class="col">
                        <q-btn square label="SHA512" color="primary" @click="encodeSHA512" class="full-width" />
                    </div>
                    <div class="col">
                        <q-btn square label="MD5" color="primary" @click="encodeMD5" class="full-width" />
                    </div>
                    <div class="col">
                        <q-btn square label="复制HASH结果" color="purple" @click="copyHASHOutput" :disable="!hashOutputText" class="full-width" />
                    </div>
                    <div class="col">
                        <q-btn square label="下载为文件" color="amber" @click="downloadHASHOutput" :disable="!hashOutputText" class="full-width" />
                    </div>
                    <div class="col">
                        <q-btn square label="清除" color="negative" @click="clearHASH" class="full-width" />
                    </div>
                </div>
            </q-card-section>

            <q-separator inset />

            <q-card-section class="q-pt-none">
                <div class="row q-gutter-md q-mt-sm">
                    <div class="col">
                        <q-input v-model="hashInputText" type="textarea" autogrow outlined square label="输入明文" :placeholder="placeholder" ref="txtHashInputText" />
                    </div>
                    <div class="col">
                        <q-input v-model="hashOutputText" type="textarea" autogrow outlined square label="HASH结果" readonly :placeholder="placeholder" />
                    </div>
                </div>
            </q-card-section>
        </q-card>

        <q-card flat square bordered class="my-card q-mt-md">
            <q-card-section>
                <div class="text-h4 text-deep-orange">Native → ASCII</div>
            </q-card-section>

            <q-card-section class="q-pt-none">
                <div class="row q-gutter-md">
                    <div class="col">
                        <div class="row q-gutter-md">
                            <div class="col">
                                <q-btn square label="Native → ASCII" class="full-width" @click="nativeToAscii" color="primary"></q-btn>
                            </div>
                            <div class="col">
                                <q-btn square label="清空" class="full-width" @click="clearNativeToAsciiInput" color="negative"></q-btn>
                            </div>
                            <div class="col">
                                <q-btn square label="复制" class="full-width" @click="copyNativeToAsciiInput" color="purple"></q-btn>
                            </div>
                        </div>

                        <div class="row q-gutter-md q-mt-sm">
                            <div class="col">
                                <q-input v-model="nativeToAsciiInput" type="textarea" autogrow outlined square label="输入Native" ref="txtNativeToAsciiInput" />
                            </div>
                        </div>
                    </div>

                    <div class="col">
                        <div class="row q-gutter-md">
                            <div class="col">
                                <q-btn square label="ASCII → Native" class="full-width" @click="asciiToNative" color="primary"></q-btn>
                            </div>
                            <div class="col">
                                <q-btn square label="清空" class="full-width" @click="clearNativeToAsciiOutput" color="negative"></q-btn>
                            </div>
                            <div class="col">
                                <q-btn square label="复制" class="full-width" @click="copyNativeToAsciiOutput" color="purple"></q-btn>
                            </div>
                        </div>

                        <div class="row q-gutter-md q-mt-sm">
                            <div class="col">
                                <q-input v-model="nativeToAsciiOutput" type="textarea" autogrow outlined square label="输入ASCII" ref="txtNativeToAsciiOutput" />
                            </div>
                        </div>
                    </div>
                </div>
            </q-card-section>
        </q-card>

        <q-card flat square bordered class="my-card q-mt-md">
            <q-card-section>
                <div class="text-h4 text-deep-orange">Native → UTF-8</div>
            </q-card-section>

            <q-card-section class="q-pt-none">
                <div class="row q-gutter-md">
                    <div class="col">
                        <div class="row q-gutter-md">
                            <div class="col">
                                <q-btn square label="Native → UTF-8(HTML)" class="full-width" @click="nativeToUTF8HTML" color="primary"></q-btn>
                            </div>
                            <div class="col">
                                <q-btn square label="Native → UTF-8(HEX)" class="full-width" @click="nativeToUTF8HEX" color="primary"></q-btn>
                            </div>
                            <div class="col">
                                <q-btn square label="清空" class="full-width" @click="clearNativeToUTF8Input" color="negative"></q-btn>
                            </div>
                            <div class="col">
                                <q-btn square label="复制" class="full-width" @click="copyNativeToUTF8Input" color="purple"></q-btn>
                            </div>
                        </div>
                        <div class="row q-gutter-md q-mt-sm">
                            <div class="col">
                                <q-input v-model="nativeToUTF8Input" type="textarea" autogrow outlined square label="输入Native" ref="txtNativeToUTF8Input" />
                            </div>
                        </div>
                    </div>

                    <div class="col">
                        <div class="row q-gutter-md">
                            <div class="col">
                                <q-btn square label="UTF-8 → Native(HTML)" class="full-width" @click="utf8ToNativeHTML" color="primary"></q-btn>
                            </div>
                            <div class="col">
                                <q-btn square label="UTF-8 → Native(HEX)" class="full-width" @click="utf8ToNativeHEX" color="primary"></q-btn>
                            </div>
                            <div class="col">
                                <q-btn square label="清空" class="full-width" @click="clearNativeToUTF8Output" color="negative"></q-btn>
                            </div>
                            <div class="col">
                                <q-btn square label="复制" class="full-width" @click="copyNativeToUTF8Output" color="purple"></q-btn>
                            </div>
                        </div>
                        <div class="row q-gutter-md q-mt-sm">
                            <div class="col">
                                <q-input v-model="nativeToUTF8Output" type="textarea" autogrow outlined square label="输入UTF-8" ref="txtNativeToUTF8Output" />
                            </div>
                        </div>
                    </div>
                </div>
            </q-card-section>
        </q-card>

        <q-card flat square bordered class="my-card q-mt-md">
            <q-card-section>
                <div class="text-h4 text-deep-orange">Native → Unicode</div>
            </q-card-section>

            <q-card-section class="q-pt-none">
                <div class="row q-gutter-md">
                    <div class="col">
                        <div class="row q-gutter-md">
                            <div class="col">
                                <q-btn square label="Native → Unicode(HTML)" class="full-width" @click="nativeToUnicodeHTML" color="primary"></q-btn>
                            </div>
                            <div class="col">
                                <q-btn square label="Native → Unicode(HEX)" class="full-width" @click="nativeToUnicodeHEX" color="primary"></q-btn>
                            </div>
                            <div class="col">
                                <q-btn square label="清空" class="full-width" @click="clearNativeToUnicodeInput" color="negative"></q-btn>
                            </div>
                            <div class="col">
                                <q-btn square label="复制" class="full-width" @click="copyNativeToUnicodeInput" color="purple"></q-btn>
                            </div>
                        </div>
                        <div class="row q-gutter-md q-mt-sm">
                            <div class="col">
                                <q-input v-model="nativeToUnicodeInput" type="textarea" autogrow outlined square label="输入Native" ref="txtNativeToUnicodeInput" />
                            </div>
                        </div>
                    </div>

                    <div class="col">
                        <div class="row q-gutter-md">
                            <div class="col">
                                <q-btn square label="Unicode → Native(HTML)" class="full-width" @click="unicodeToNativeHTML" color="primary"></q-btn>
                            </div>
                            <div class="col">
                                <q-btn square label="Unicode → Native(HEX)" class="full-width" @click="unicodeToNativeHEX" color="primary"></q-btn>
                            </div>
                            <div class="col">
                                <q-btn square label="清空" class="full-width" @click="clearNativeToUnicodeOutput" color="negative"></q-btn>
                            </div>
                            <div class="col">
                                <q-btn square label="复制" class="full-width" @click="copyNativeToUnicodeOutput" color="purple"></q-btn>
                            </div>
                        </div>
                        <div class="row q-gutter-md q-mt-sm">
                            <div class="col">
                                <q-input v-model="nativeToUnicodeOutput" type="textarea" autogrow outlined square label="输入UTF-8" ref="txtNativeToUnicodeOutput" />
                            </div>
                        </div>
                    </div>
                </div>
            </q-card-section>
        </q-card>
    </q-page>
</template>

<script setup>
import { ref } from "vue";
import notify from "src/utils/notify";
import { sha224, sha256 } from "js-sha256";
import { sha384, sha512 } from "js-sha512";
import md5 from "js-md5";
import clipboard from "src/utils/clipboard";
import download from "src/utils/download";

const base64InputText = ref("");
const base64OutputText = ref("");
const hashInputText = ref("");
const hashOutputText = ref("");
const placeholder = '输入普通文本然后点击 "文本 → Base64"，或粘贴 Base64 再点击 "Base64 → 文本"';
const nativeToAsciiInput = ref("");
const nativeToAsciiOutput = ref("");
const nativeToUTF8Input = ref("");
const nativeToUTF8Output = ref("");
const nativeToUnicodeInput = ref("");
const nativeToUnicodeOutput = ref("");

const txtBase64InputText = ref(null);
const txtHashInputText = ref(null);
const txtNativeToAsciiInput = ref(null);
const txtNativeToAsciiOutput = ref(null);
const txtNativeToUTF8Input = ref(null);
const txtNativeToUTF8Output = ref(null);
const txtNativeToUnicodeInput = ref(null);
const txtNativeToUnicodeOutput = ref(null);

function utf8ToBase64(str) {
    const bytes = new TextEncoder().encode(str);
    let binary = "";
    for (let i = 0; i < bytes.length; i++) binary += String.fromCharCode(bytes[i]);
    return btoa(binary);
}

function base64ToUtf8(b64) {
    const binary = atob(b64);
    const len = binary.length;
    const bytes = new Uint8Array(len);
    for (let i = 0; i < len; i++) bytes[i] = binary.charCodeAt(i);
    return new TextDecoder().decode(bytes);
}

function normalizeBase64(s) {
    return s.replace(/\s+/g, "");
}

function isLikelyBase64(s) {
    const n = normalizeBase64(s);
    return n.length > 0 && /^[A-Za-z0-9+/]+={0,2}$/.test(n);
}

function encodeText() {
    try {
        if (!base64InputText.value) {
            notify.error("请输入要编码的文本");
            return;
        }
        base64OutputText.value = utf8ToBase64(base64InputText.value);
        notify.ok("编码完成");
    } catch (e) {
        notify.error(`编码失败：${e.message || e}`);
    }
}

function decodeText() {
    try {
        if (!base64InputText.value) {
            notify.error("请输入 Base64 内容");
            return;
        }

        // 去掉 data:*;base64, 前缀（如果存在），并去除空白
        let s = base64InputText.value.trim();
        const m = s.match(/^data:[^;]+;base64,(.+)$/i);
        if (m) s = m[1];

        const candidate = normalizeBase64(s);
        if (!isLikelyBase64(candidate)) {
            notify.error("输入看起来不像 Base64 字符串");
            return;
        }
        base64OutputText.value = base64ToUtf8(candidate);
        notify.ok("解码完成");
    } catch (e) {
        notify.error(`解码失败：${e.message || e}`);
    }
}

function copyBase64Output() {
    try {
        clipboard.copyToClipboard(base64OutputText.value);
        notify.ok("Base64结果已复制到剪贴板");
    } catch (e) {
        notify.error(e.message || e);
        return;
    }
}

const downloadBase64Output = () => {
    try {
        download.saveToFile(base64OutputText.value, "base64.txt");
    } catch (e) {
        notify.error(e.message || e);
        return;
    }
};

function clearAll() {
    base64InputText.value = "";
    base64OutputText.value = "";
    txtBase64InputText.value?.focus();
}

const encodeSHA1 = async () => {
    try {
        if (!hashInputText.value) {
            notify.error("请输入要计算哈希的文本");
            return;
        }
        const data = new TextEncoder().encode(hashInputText.value);
        const hashBuffer = await crypto.subtle.digest("SHA-1", data);
        hashOutputText.value = bufferToHex(hashBuffer);
        notify.ok("SHA-1 计算完成");
    } catch (e) {
        notify.error(`SHA-1 计算失败：${e.message || e}`);
    }
};

const bufferToHex = buffer =>
    Array.from(new Uint8Array(buffer))
        .map(b => b.toString(16).padStart(2, "0"))
        .join("");

const encodeSHA224 = () => {
    try {
        if (!hashInputText.value) {
            notify.error("请输入要计算哈希的文本");
            return;
        }
        // js-sha256 的 sha224 接受字符串并返回 hex
        hashOutputText.value = sha224(hashInputText.value);
        notify.ok("SHA-224 计算完成");
    } catch (e) {
        notify.error(`SHA-224 计算失败：${e.message || e}`);
    }
};

const encodeSHA256 = () => {
    try {
        if (!hashInputText.value) {
            notify.error("请输入要计算哈希的文本");
            return;
        }
        // 使用 js-sha256 的 sha256，返回 hex 字符串
        hashOutputText.value = sha256(hashInputText.value);
        notify.ok("SHA-256 计算完成");
    } catch (e) {
        notify.error(`SHA-256 计算失败：${e.message || e}`);
    }
};

const encodeSHA384 = () => {
    try {
        if (!hashInputText.value) {
            notify.error("请输入要计算哈希的文本");
            return;
        }
        hashOutputText.value = sha384(hashInputText.value);
        notify.ok("SHA-384 计算完成");
    } catch (e) {
        notify.error(`SHA-384 计算失败：${e.message || e}`);
    }
};

const encodeSHA512 = () => {
    try {
        if (!hashInputText.value) {
            notify.error("请输入要计算哈希的文本");
            return;
        }
        hashOutputText.value = sha512(hashInputText.value);
        notify.ok("SHA-512 计算完成");
    } catch (e) {
        notify.error(`SHA-512 计算失败：${e.message || e}`);
    }
};

const encodeMD5 = () => {
    try {
        if (!hashInputText.value) {
            notify.error("请输入要计算哈希的文本");
            return;
        }
        hashOutputText.value = md5(hashInputText.value);
        notify.ok("MD5 计算完成");
    } catch (e) {
        notify.error(`MD5 计算失败：${e.message || e}`);
    }
};

const copyHASHOutput = () => {
    try {
        clipboard.copyToClipboard(hashOutputText.value);
        notify.ok("HASH结果已复制到剪贴板");
    } catch (e) {
        notify.error(e.message || e);
        return;
    }
};

const downloadHASHOutput = () => {
    try {
        download.saveToFile(hashOutputText.value, "hash.txt");
    } catch (e) {
        notify.error(e.message || e);
        return;
    }
};

const clearHASH = () => {
    hashInputText.value = "";
    hashOutputText.value = "";
    txtHashInputText.value?.focus();
};

const nativeToAscii = () => {
    try {
        if (!nativeToAsciiInput.value) {
            notify.error("请输入要转换的文本");
            return;
        }
        let result = "";
        for (let i = 0; i < nativeToAsciiInput.value.length; i++) {
            const charCode = nativeToAsciiInput.value.charCodeAt(i);
            if (charCode > 127) {
                result += "\\u" + charCode.toString(16).padStart(4, "0");
            } else {
                result += nativeToAsciiInput.value.charAt(i);
            }
        }
        nativeToAsciiOutput.value = result;
        notify.ok("转换完成");
    } catch (e) {
        notify.error(`转换失败：${e.message || e}`);
    }
};

const asciiToNative = () => {
    if (!nativeToAsciiOutput.value) {
        notify.error("请输入要转换的文本");
        return;
    }
    nativeToAsciiInput.value = nativeToAsciiOutput.value.replace(/\\u([\dA-Fa-f]{4})/g, (match, grp) => {
        return String.fromCharCode(parseInt(grp, 16));
    });
    notify.ok("转换完成");
};

const clearNativeToAsciiInput = () => {
    nativeToAsciiInput.value = "";
    txtNativeToAsciiInput.value?.focus();
};
const clearNativeToAsciiOutput = () => {
    nativeToAsciiOutput.value = "";
    txtNativeToAsciiOutput.value?.focus();
};

const copyNativeToAsciiInput = () => {
    try {
        clipboard.copyToClipboard(nativeToAsciiInput.value);
        notify.ok("输入内容已复制到剪贴板");
    } catch (e) {
        notify.error(e.message || e);
        return;
    }
};

const copyNativeToAsciiOutput = () => {
    try {
        clipboard.copyToClipboard(nativeToAsciiOutput.value);
        notify.ok("输出内容已复制到剪贴板");
    } catch (e) {
        notify.error(e.message || e);
        return;
    }
};

const nativeToUTF8HTML = () => {
    try {
        if (!nativeToUTF8Input.value) {
            notify.error("请输入要转换的文本");
            return;
        }
        nativeToUTF8Output.value = Array.from(nativeToUTF8Input.value)
            .map(ch => `&#x${ch.codePointAt(0).toString(16).toUpperCase()};`)
            .join("");
        notify.ok("转换为 UTF-8(HTML) 完成");
    } catch (e) {
        notify.error(`转换失败：${e.message || e}`);
    }
};

const nativeToUTF8HEX = () => {
    try {
        if (!nativeToUTF8Input.value) {
            notify.error("请输入要转换的文本");
            return;
        }
        nativeToUTF8Output.value = Array.from(nativeToUTF8Input.value)
            .map(ch => ch.codePointAt(0).toString(16).toUpperCase().padStart(4, "0"))
            .join(" ");
        notify.ok("转换为 UTF-8(HEX) 完成");
    } catch (e) {
        notify.error(`转换失败：${e.message || e}`);
    }
};

const utf8ToNativeHTML = () => {
    try {
        if (!nativeToUTF8Output.value) {
            notify.error("请输入要转换的文本");
            return;
        }
        nativeToUTF8Input.value = nativeToUTF8Output.value.replace(/&#x([\dA-Fa-f]+);/g, (match, grp) => {
            return String.fromCodePoint(parseInt(grp, 16));
        });
        notify.ok("转换为 Native 完成");
    } catch (e) {
        notify.error(`转换失败：${e.message || e}`);
    }
};

const utf8ToNativeHEX = () => {
    try {
        if (!nativeToUTF8Output.value) {
            notify.error("请输入要转换的文本");
            return;
        }
        nativeToUTF8Input.value = nativeToUTF8Output.value
            .split(/\s+/)
            .map(h => String.fromCodePoint(parseInt(h, 16)))
            .join("");
        notify.ok("转换为 Native 完成");
    } catch (e) {
        notify.error(`转换失败：${e.message || e}`);
    }
};

const clearNativeToUTF8Input = () => {
    nativeToUTF8Input.value = "";
    txtNativeToUTF8Input.value?.focus();
};
const clearNativeToUTF8Output = () => {
    nativeToUTF8Output.value = "";
    txtNativeToUTF8Output.value?.focus();
};

const copyNativeToUTF8Input = () => {
    try {
        clipboard.copyToClipboard(nativeToUTF8Input.value);
        notify.ok("输入内容已复制到剪贴板");
    } catch (e) {
        notify.error(e.message || e);
        return;
    }
};

const copyNativeToUTF8Output = () => {
    try {
        clipboard.copyToClipboard(nativeToUTF8Output.value);
        notify.ok("输出内容已复制到剪贴板");
    } catch (e) {
        notify.error(e.message || e);
        return;
    }
};

const nativeToUnicodeHTML = () => {
    if (!nativeToUnicodeInput.value) {
        notify.error("请输入要转换的文本");
        return;
    }

    nativeToUnicodeOutput.value = Array.from(nativeToUnicodeInput.value)
        .map(char => {
            return `&#${char.codePointAt(0)};`;
        })
        .join("");

    notify.ok("转换为 Unicode(HTML) 完成");
};

const nativeToUnicodeHEX = () => {
    if (!nativeToUnicodeInput.value) {
        notify.error("请输入要转换的文本");
        return;
    }

    nativeToUnicodeOutput.value = Array.from(nativeToUnicodeInput.value)
        .map(char => {
            return char.codePointAt(0).toString(16).toUpperCase().padStart(4, "0");
        })
        .join(" ");

    notify.ok("转换为 Unicode(HEX) 完成");
};

const unicodeToNativeHTML = () => {
    if (!nativeToUnicodeOutput.value) {
        notify.error("请输入要转换的文本");
        return;
    }

    nativeToUnicodeInput.value = nativeToUnicodeOutput.value.replace(/&#(\d+);/g, (match, grp) => {
        return String.fromCodePoint(parseInt(grp, 10));
    });

    notify.ok("转换为 Native 完成");
};

const unicodeToNativeHEX = () => {
    if (!nativeToUnicodeOutput.value) {
        notify.error("请输入要转换的文本");
        return;
    }

    nativeToUnicodeInput.value = nativeToUnicodeOutput.value
        .split(/\s+/)
        .map(h => String.fromCodePoint(parseInt(h, 16)))
        .join("");

    notify.ok("转换为 Native 完成");
};

const clearNativeToUnicodeInput = () => {
    nativeToUnicodeInput.value = "";
    txtNativeToUnicodeInput.value?.focus();
};
const clearNativeToUnicodeOutput = () => {
    nativeToUnicodeOutput.value = "";
    txtNativeToUnicodeOutput.value?.focus();
};

const copyNativeToUnicodeInput = () => {
    try {
        clipboard.copyToClipboard(nativeToUnicodeInput.value);
        notify.ok("输入内容已复制到剪贴板");
    } catch (e) {
        notify.error(e.message || e);
        return;
    }
};

const copyNativeToUnicodeOutput = () => {
    try {
        clipboard.copyToClipboard(nativeToUnicodeOutput.value);
        notify.ok("输出内容已复制到剪贴板");
    } catch (e) {
        notify.error(e.message || e);
        return;
    }
};
</script>

<style scoped>
/* 简单样式，可按需微调 */
</style>
