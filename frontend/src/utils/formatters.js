export function formatD(date) {
    var d = new Date(date),
        month = "" + (d.getMonth() + 1),
        day = "" + d.getDate(),
        year = d.getFullYear();

    if (month < 10) month = "0" + month;
    if (day < 10) day = "0" + day;

    return year + "-" + month + "-" + day;
}

export function formatT(date) {
    var d = new Date(date),
        h = d.getHours(),
        m = d.getMinutes(),
        s = d.getSeconds(),
        ms = d.getMilliseconds();

    if (h < 10) h = "0" + h;
    if (m < 10) m = "0" + m;
    if (s < 10) s = "0" + s;
    if (ms < 10) {
        ms = "00" + ms;
    }
    if (ms >= 10 && ms < 100) {
        ms = "0" + ms;
    }

    return h + ":" + m + ":" + s + "." + ms;
}