import { BASEPATH } from "@/api/core.js"

export default {
    getHolidays() {
        let res = {
            today: "",
            lastDay: "",
            isUpdateRequired: false,
            isHolidayToday: false,
            dates: [],
            userNotifyRequired: false,
        };
        fetch(`${BASEPATH}/holidays`)
            .then((r) => r.json())
            .then(function (j) {
                res.today = j.today;
                res.lastDay = "";
                res.isUpdateRequired = j.isUpdateRequired;
                res.isHolidayToday = j.isHolidayToday;

                if (j.dates && j.dates.length > 0) {
                    res.dates.push(...j.dates);
                    res.lastDay = j.dates[j.dates.length - 1].dt;
                    res.userNotifyRequired = res.isHolidayToday;// && !window.webpackHotUpdate;
                }
            })
        return res;
    }
};