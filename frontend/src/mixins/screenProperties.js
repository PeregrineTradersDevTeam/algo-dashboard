
const screenProperties = {
    get() {
        let res = {
            width: window.innerWidth,
            height: window.innerHeight,
            type: "desktop",
            mode: "landscape",
        };

        if (window.orientation == 90 || window.orientation == -90) {
            res.mode = "landscape";
            res.type = "mobile";
        } else {
            if (window.matchMedia("only screen and (max-width: 492px)").matches) {
                res.mode = "portrait";
                res.type = "mobile";
            }
        }
        return res;
    },
    methods: {
        isMobile() {
            return screenProperties.get().type === "mobile"
        },
        isPortrait() {
            return screenProperties.get().mode === "portrait"
        }
    }
}

export default screenProperties;