Vue.component('gosha-nav', {
    props:['page'],
    template: `
        <div class="nav" v-on:click="navclick">
            <div class="nav-page-icon mdi" v-bind:class="iconclass"></div>
        </div>
    `,
    methods: {
        navclick: function() {
            this.$root.currentPage = this.page.id
            this.$root.screenclick()
        }
    },
    computed: {
        iconclass: function() {
            var iconobj = {}
            if (this.page.icon) {
                iconobj[this.page.icon] = true
            } else {
                iconobj['mdi-flask-empty-remove-outline'] = true
            }
            return iconobj
        }
    }
})

Vue.component('gosha-page', {
    props:['page'],
    template: `
        <div class="page">
            <div class="page-title">{{ page.title }}</div>
            <slot></slot>
        </div>
    `
})

Vue.component('gosha-group', {
    props:['group'],
    template: `
        <div class="group">
            <slot></slot>
        </div>
    `
})

Vue.component('gosha-gauge', {
    props: ['card'],
    template: `
        <div class="card gauge-card off">
            <vue-svg-gauge
                v-bind:value="card.state"
                v-bind:min="card.min"
                v-bind:max="card.max"
                :start-angle="-160" 
                :end-angle=160
                :separator-step="0"
                :separator-thickness="1"
                :scale-interval="0"
                :inner-radius="80"
                base-color="#074133">
                </vue-svg-gauge>
            <div class="sensor-value">{{ card.state }} {{ card.unit_of_measurement }}</div>
            <div class="card-title">{{ card.title }}</div>
        </div>
    `,
    methods: {
        datacall: function(evt) {
            if (evt.type && evt.type == "event" && evt.event && evt.event.data &&
                evt.event.data.entity_id == this.card.id && evt.event.data.new_state) {
                    this.card.state = evt.event.data.new_state.state
            }

            if (evt.entity_id && evt.entity_id == this.card.id) {
                this.card.state = evt.state
            }
        }
    },
    mounted: function() {
        this.$root.datacalls[this.card.id] = this.datacall

        if (this.$root.lastevents[this.card.id]) {
            this.datacall(this.$root.lastevents[this.card.id])
        }
    }
})


Vue.component('gosha-switch', {
    props: ['card'],
    template: `
        <div class="card" v-bind:class="stateclass" v-on:click="cardclick">
            <div class="card-icon mdi" v-bind:class="iconclass"></div>
            <div class="card-title">{{ card.title }}</div>
        </div>
    `,
    data: function() {
        return {
            msgid: 0
        }
    },
    computed: {
        stateclass: function() {
            return {
                'on': this.card.state == 'on',
                'off': this.card.state != 'on'
            }
        },
        iconclass: function() {
            var iconobj = {}
            if (this.card.icon) {
                iconobj[this.card.icon] = true
            } else {
                iconobj['mdi-flask-empty-remove-outline'] = true
            }
            return iconobj
        }
    },
    methods: {
        cardclick: function() {
            this.$root.screenclick()

            if (this.$root.wsconn) {
                this.msgid++
                var cmd = {
                    "type": "call_service",
                    "domain": this.card.domain,
                    "service": "toggle",
                    "service_data":
                        {
                            "entity_id": this.card.id
                        },
                        "id": this.msgid
                }
                this.$root.wsconn.send(JSON.stringify(cmd))
            }
        },
        datacall: function(evt) {
            if (evt.type && evt.type == "event" && evt.event && evt.event.data &&
                evt.event.data.entity_id == this.card.id && evt.event.data.new_state) {
                    this.card.state = evt.event.data.new_state.state
            }

            if (evt.entity_id && evt.entity_id == this.card.id) {
                this.card.state = evt.state
            }
        }
    },
    mounted: function() {
        this.$root.datacalls[this.card.id] = this.datacall

        if (this.$root.lastevents[this.card.id]) {
            this.datacall(this.$root.lastevents[this.card.id])
        }
    }
})

Vue.component('gosha-light', {
    props: ['card'],
    template: `
    <div class="card" v-bind:class="stateclass" v-on:click="cardclick">
        <div class="card-icon mdi" v-bind:class="iconclass"></div>
        <div class="card-title">{{ card.title }}</div>
        <div class="card-brightness-container" v-show="card.state == 'on'">
            <div class="brightness-control">{{card.brightness}}</div>
            <div class="brightness-control up" v-on:click="upclick" v-show="card.brightness < card.max">+</div>
            <div class="brightness-control down" v-on:click="downclick" v-show="card.brightness > card.min">-</div>
        </div>
    </div>
    `,
    data : function() {
        return {
            msgid: 0
        }
    },
    computed: {
        stateclass: function() {
            return {
                'on': this.card.state == 'on',
                'off': this.card.state != 'on'
            }
        },
        iconclass: function() {
            var iconobj = {}
            if (this.card.icon) {
                iconobj[this.card.icon] = true
            } else {
                iconobj['mdi-flask-empty-remove-outline'] = true
            }
            return iconobj
        }
    },
    methods: {
        cardclick: function() {
            this.$root.screenclick()

            if (this.$root.wsconn) {
                this.msgid++
                var cmd = {
                    "type": "call_service",
                    "domain": this.card.domain,
                    "service": "toggle",
                    "service_data":
                        {
                            "entity_id": this.card.id
                        },
                        "id": this.msgid
                }
                this.$root.wsconn.send(JSON.stringify(cmd))
            }
        },
        upclick: function(event) {
            event.stopPropagation()
            brightness = this.card.brightness + this.card.step
            if (brightness > this.card.max) {
                brightness = this.card.max
            }

            if (this.$root.wsconn) {
                this.msgid++
                var cmd = {
                    "type": "call_service",
                    "domain": this.card.domain,
                    "service": '' + brightness,
                    "service_data":
                        {
                            "entity_id": this.card.id
                        },
                        "id": this.msgid
                }
                this.$root.wsconn.send(JSON.stringify(cmd))
            }
        },
        downclick: function(event) {
            event.stopPropagation()
            brightness = this.card.brightness - this.card.step
            if (brightness < this.card.min) {
                brightness = this.card.min
            }
            
            if (this.$root.wsconn) {
                this.msgid++
                var cmd = {
                    "type": "call_service",
                    "domain": this.card.domain,
                    "service": '' + brightness,
                    "service_data":
                        {
                            "entity_id": this.card.id
                        },
                        "id": this.msgid
                }
                this.$root.wsconn.send(JSON.stringify(cmd))
            }
        },
        datacall: function(evt) {
            if (evt.type && evt.type == "event" && evt.event && evt.event.data &&
                evt.event.data.entity_id == this.card.id && evt.event.data.new_state) {
                    this.card.state = evt.event.data.new_state.state
                    this.card.brightness = evt.event.data.new_state.attributes['brightness']
            }

            if (evt.entity_id && evt.entity_id == this.card.id) {
                this.card.state = evt.state
                this.card.brightness = evt.attributes['brightness']
            }
        }
    },
    mounted: function() {
        this.$root.datacalls[this.card.id] = this.datacall

        if (this.$root.lastevents[this.card.id]) {
            this.datacall(this.$root.lastevents[this.card.id])
        }
    }
})

Vue.component('gosha-weather', {
    props: ['card'],
    template: `
        <div class="card weather-card off">
            <div class="weather-state">{{ card.state }}</div>
            <div class="weather-icon-container">
                <i class="weather-icon owi" v-bind:class="iconclass"></i>
            </div>
            <div class="weather-temp">{{ card.attributes.temp }} ℃</div>
            <div class="weather-attrs"><i class="mdi mdi-water-percent"></i> {{ card.attributes.humidity }}%   <i class="mdi mdi-weather-windy"></i> {{ card.attributes.wind_speed }}км/ч</div>
        </div>
    `,
    computed: {
        iconclass: function() {
            var iconobj = {}
            if (this.card.attributes.icon) {
                iconobj['owi-' + this.card.attributes.icon] = true
            } else {
                iconobj['owi-01d'] = true
            }
            return iconobj
        }
    },
    methods: {
        datacall: function(evt) {
            if (evt.type && evt.type == "event" && evt.event && evt.event.data &&
                evt.event.data.entity_id == this.card.id && evt.event.data.new_state) {
                    this.card.state = evt.event.data.new_state.state
                    this.card.attributes = evt.event.data.new_state.attributes
            }

            if (evt.entity_id && evt.entity_id == this.card.id) {
                this.card.state = evt.state
                this.card.attributes = evt.attributes
            }
        }
    },
    mounted: function() {
        this.$root.datacalls[this.card.id] = this.datacall

        if (this.$root.lastevents[this.card.id]) {
            this.datacall(this.$root.lastevents[this.card.id])
        }
    }
})

Vue.component('gosha-service', {
    props: ['card'],
    template: `
        <div class="card gauge-card" v-bind:class="stateclass" v-on:click="cardclick">
            <div class="card-icon mdi" v-bind:class="iconclass"></div>
            <div class="card-title">{{ card.title }}</div>
        </div>
    `,
    data: function() {
        return {
            msgid: 0
        }
    },
    computed: {
        iconclass: function() {
            var iconobj = {}
            if (this.card.icon) {
                iconobj[this.card.icon] = true
            } else {
                iconobj['mdi-flask-empty-remove-outline'] = true
            }
            return iconobj
        },
        stateclass: function() {
            return {
                'on': this.card.state == 'on',
                'off': this.card.state != 'on'
            }
        }
    },
    methods: {
        cardclick: function() {
            this.$root.screenclick()

            if (this.$root.wsconn) {
                this.msgid++
                var cmd = {
                    "type": "call_service",
                    "domain": this.card.service,
                    "service": this.card.method,
                    "service_data":
                        {
                            "entity_id": this.card.id
                        },
                        "id": this.msgid
                }
                this.$root.wsconn.send(JSON.stringify(cmd))

                this.card.state = 'on'
                setTimeout(function(o) {
                    o.card.state = 'off'
                }, 500, this)
            }
        }
    }
})

Vue.component('gosha-alarm', {
    props: ['card'],
    template: `
    <div class="card off" v-on:click="cardclick">
        <div class="card-icon mdi" v-bind:class="iconclass"></div>
        <div class="card-title">{{ card.title }}</div>
        <div class="alarm-mode-popup" v-if="popup">
            <div class="card off" v-on:click="popclick('disarm', $event)">
                <div class="card-icon mdi mdi-bell-off-outline"></div>
                <div class="card-title">Отключено</div>
            </div>
            <div class="card off" v-on:click="popclick('arm_home', $event)">
                <div class="card-icon mdi mdi-bell-outline"></div>
                <div class="card-title">Дома</div>
            </div>
            <div class="card off" v-on:click="popclick('arm_away', $event)">
                <div class="card-icon mdi mdi-bell-check-outline"></div>
                <div class="card-title">Ушли</div>
            </div>
            <div class="card off" v-on:click="popclick('arm_night', $event)">
                <div class="card-icon mdi mdi-bell-sleep-outline"></div>
                <div class="card-title">Ночь</div>
            </div>
        </div>
    </div>
    `,
    data: function() {
        return {
            msgid: 0,
            popup: false,
        }
    },
    computed: {
        iconclass: function() {
            var iconobj = {}
            switch (this.card.state) {
                case 'disarmed':
                    iconobj['mdi-bell-off-outline'] = true
                    break
                case 'armed_home':
                    iconobj['mdi-bell-outline'] = true
                    break
                case 'armed_away':
                    iconobj['mdi-bell-check-outline'] = true
                    break
                case 'armed_night':
                    iconobj['mdi-bell-sleep-outline'] = true
                    break
                case 'triggered':
                    iconobj['mdi-bell-alert-outline'] = true
                    break
                default:
                    iconobj[card.icon] = true
            }
            return iconobj
        }
    },
    methods: {
        cardclick: function() {
            this.$root.screenclick()
            this.popup = !this.popup
        },
        popclick: function(state, event) {
            event.stopPropagation()

            if (this.$root.wsconn) {
                this.msgid++
                var cmd = {
                    type: "call_service",
                    domain: this.card.domain,
                    service: state,
                    service_data: {
                        entity_id: this.card.id
                    },
                    id: this.msgid
                }

                this.$root.wsconn.send(JSON.stringify(cmd))
            }

            this.popup = !this.popup
        },
        datacall: function(evt) {
            if (evt.type && evt.type == "event" && evt.event && evt.event.data &&
                evt.event.data.entity_id == this.card.id && evt.event.data.new_state) {
                    this.card.state = evt.event.data.new_state.state
            }

            if (evt.entity_id && evt.entity_id == this.card.id) {
                this.card.state = evt.state
            }
        },
    },
    mounted: function() {
            this.$root.datacalls[this.card.id] = this.datacall

            if (this.$root.lastevents[this.card.id]) {
                this.datacall(this.$root.lastevents[this.card.id])
            }
        }
})

Vue.component('gosha-sensor', {
    props: ['card'],
    template: `
        <div class="card off">
            <div class="sensor-state">{{ card.state }} {{ card.attributes.unit_of_measurement }}</div>
            <div class="card-title">{{ card.title }}</div>
        </div>
    `,
    methods: {
        datacall: function(evt) {
            if (evt.type && evt.type == "event" && evt.event && evt.event.data &&
                evt.event.data.entity_id == this.card.id && evt.event.data.new_state) {
                    this.card.state = evt.event.data.new_state.state
                    this.card.attributes = evt.event.data.new_state.attributes
            }

            if (evt.entity_id && evt.entity_id == this.card.id) {
                this.card.state = evt.state
                this.card.attributes = evt.attributes
            }
        }
    },
    mounted: function() {
        this.$root.datacalls[this.card.id] = this.datacall

        if (this.$root.lastevents[this.card.id]) {
            this.datacall(this.$root.lastevents[this.card.id])
        }
    }
})

Vue.component('gosha-binary-sensor', {
    props: ['card'],
    template: `
        <div class="card" v-bind:class="stateclass">
            <div class="card-icon mdi" v-bind:class="iconclass"></div>
            <div class="card-title">{{ card.title }}</div>
        </div>
    `,
    computed: {
        stateclass: function() {
            return {
                'on': this.card.state == 'on',
                'off': this.card.state != 'on'
            }
        },
        iconclass: function() {
            var iconobj = {}
            if (this.card.icon) {
                iconobj[this.card.icon] = true
            } else {
                iconobj['mdi-flask-empty-remove-outline'] = true
            }
            return iconobj
        }
    },
    methods: {
        datacall: function(evt) {
            if (evt.type && evt.type == "event" && evt.event && evt.event.data &&
                evt.event.data.entity_id == this.card.id && evt.event.data.new_state) {
                    this.card.state = evt.event.data.new_state.state
            }

            if (evt.entity_id && evt.entity_id == this.card.id) {
                this.card.state = evt.state
            }
        }
    },
    mounted: function() {
        this.$root.datacalls[this.card.id] = this.datacall

        if (this.$root.lastevents[this.card.id]) {
            this.datacall(this.$root.lastevents[this.card.id])
        }
    }
})

var $this

window.addEventListener("load", function(event) {

new Vue({
    el: '#gosha',
    data: {
        pages: CONFIG.pages,
        currentPage: 1,
        currentTime: "21:43",
        currentDate: "",
        tickCount: 0,
        screensaverActive: false,
        datacalls: {},
        lastevents: {},
        wsconn: null,
        connected: false,
        auth: true,
        msgid: 0,
    },
    created: function() {
        console.log('Gosha is created')
    },
    methods: {
        getCurrentTime: function() {
            var d = new Date()
            var mark = d.getSeconds() % 2 == 0 ? ' : ' : "   "
            var h = (d.getHours() < 10 ? '0' : '') + d.getHours()
            var m = (d.getMinutes() < 10 ? '0' : '') + d.getMinutes()
            this.currentTime = h + mark + m
            this.currentDate = d.toLocaleDateString("ru-RU", { weekday: 'long', year: 'numeric', month: 'long', day: 'numeric' })
        },

        onconnectopen: function(evt) {
            console.log('WS Connection established')
            $this.connected = true
            $this.msgid++
            var cmd = {
                type: "get_states",
                id: $this.msgid
            }
            $this.wsconn.send(JSON.stringify(cmd))
        },

        onconnectclose: function(evt) {
            console.log('WS Connection closed')
            $this.connected = false
        },

        onsocketmessage: function(evt) {
            if (!evt.data) {
                return
            }

            var msg = JSON.parse(evt.data)
            if (!msg) {
                return
            }

            if (msg.type == 'auth_required') {
                $this.auth = false
                $this.msgid++
                var cmd = {
                    type: "auth",
                    token: CONFIG.token,
                    id: $this.msgid
                }
                $this.wsconn.send(JSON.stringify(cmd))
            }

            if (msg.type == 'auth_ok') {
                $this.auth = true

                $this.msgid++
                var cmd = {
                    type: "get_states",
                    id: $this.msgid
                }
                $this.wsconn.send(JSON.stringify(cmd))
            }

            if (msg.type == "event") {
                $this.lastevents[msg.event.data.entity_id] = msg

                if ($this.datacalls[msg.event.data.entity_id]) {
                    $this.datacalls[msg.event.data.entity_id](msg)
                }
            }

            if (msg.type == "result" && msg.result) {
                msg.result.forEach(function(result) {
                    $this.lastevents[result.entity_id] = result
                    if ($this.datacalls[result.entity_id]) {
                        $this.datacalls[result.entity_id](result)
                    }
                })
            }
        },

        connect: function() {
            if (this.wsconn) {
                this.wsconn.onopen = null
                this.wsconn.onclose = null
                this.wsconn.onmessage = null
                this.wsconn.close()
                this.wsconn = null
            }

            this.wsconn = new WebSocket(("ws://" + document.location.host + "/api/ws"))
            this.wsconn.onopen = this.onconnectopen
            this.wsconn.onclose = this.onconnectclose
            this.wsconn.onmessage = this.onsocketmessage
        },

        screenclick: function() {
            this.tickCount = 0
            this.screensaverActive = false
        }
    },
    mounted: function() {
        console.log('Mounted')

        $this = this

        this.connect()

        setInterval(function() {
            $this.getCurrentTime()

            $this.tickCount++
            if ($this.tickCount > CONFIG.screensaver.timeout) {
                $this.tickCount = 0
                $this.screensaverActive = true
            }

            if ($this.wsconn.readyState === WebSocket.CLOSED) {
                $this.connected = false
            }

            if (!$this.connected) {
                $this.connect()
            }

        }, 1000)

        setInterval(function() {
            if ($this.connected) {
                $this.msgid++
                if ($this.msgid > 1000) {
                    $this.msgid = 1
                }
                var cmd = {
                    "type": "ping",
                    "id": $this.msgid
                }
                $this.wsconn.send(JSON.stringify(cmd))
            }
        }, 10000)
    }
})

}) ;