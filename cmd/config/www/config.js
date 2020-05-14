var CONFIG = {
    token: 'gOsHa666',
    screensaver: {
        timeout: 300
    },
    pages: [
        {
            id: 1,
            title: 'Дом',
            icon: 'mdi-home-circle',
            groups: [
                {
                    id: 1,
                    title: 'group 1',
                    cards: [
                        {
                            domain: 'weather',
                            id: 'weather.theweather',
                            title: 'Погода',
                            state: 'Clouds',
                            attributes: {
                                icon: '01d',
                                temp: 8.4,
                                humidity: 90,
                                wind_speed: 4,
                                pressure: 1009
                            }
                        },
                        {
                            domain: 'spacer',
                            width: 0.5
                        },
                        {
                            domain: 'switch',
                            id: 'switch.hall_outlet2', 
                            title: 'Телевизор',
                            state: 'off',
                            icon: 'mdi-television'
                        },
                        {
                            domain: 'switch',
                            id: 'switch.bathroom_fan', 
                            title: 'Ванная',
                            state: 'off',
                            icon: 'mdi-fan'
                        },
                        {
                            domain: 'switch',
                            id: 'switch.corridor_light', 
                            title: 'Прихожая',
                            state: 'off',
                            icon: 'mdi-coach-lamp'
                        },
                    ]
                },
                {
                    id: 2,
                    title: 'group 2',
                    cards: [
                        {
                            domain: 'gauge',
                            id: 'sensor.hall_temp',
                            title: 'Темп',
                            state: 12.3,
                            unit_of_measurement: '℃',
                            min: 0,
                            max: 40
                        },
                        {
                            domain: 'spacer',
                            width: 0.5
                        },
                        {
                            domain: 'light',
                            id: 'light.cabinet_backlight',
                            title: 'Стол кабинет',
                            state: 'off',
                            min: 0,
                            max: 100,
                            step: 10,
                            brightness: 0,
                            icon: 'mdi-layers-triple-outline'
                        },
                        {
                            domain: 'switch',
                            id: 'switch.hall_sofalamp', 
                            title: 'Диван',
                            state: 'off',
                            icon: 'mdi-lamp'
                        },
                        {
                            domain: 'switch',
                            id: 'switch.corridor_passlight', 
                            title: 'Проход',
                            state: 'off',
                            icon: 'mdi-outdoor-lamp'
                        },
                    ]
                }
            ]
        },
        {
            id: 2,
            title: 'Кабинет',
            icon: 'mdi-security-network',
            groups: [
                {
                    id: 1,
                    title: 'group 1',
                    cards: [
                        {
                            domain: 'switch',
                            id: 'switch.cabinet_toplight', 
                            title: 'Верхний свет',
                            state: 'off',
                            icon: 'mdi-ceiling-light'
                        },
                        {
                            domain: 'spacer',
                            width: 0.5
                        },
                        {
                            domain: 'switch',
                            id: 'switch.cabinet_tablelamp',
                            title: 'Настольная лампа ',
                            state: 'off',
                            icon: 'mdi-desk-lamp'
                        },
                        {
                            domain: 'switch',
                            id: 'switch.cabinet_outlet1',
                            title: 'Розетка',
                            state: 'off',
                            icon: 'mdi-power-socket-eu'
                        },
                        {
                            domain: 'light',
                            id: 'light.cabinet_backlight',
                            title: 'Стол',
                            state: 'off',
                            min: 0,
                            max: 100,
                            step: 10,
                            brightness: 0,
                            icon: 'mdi-layers-triple-outline'
                        },
                    ]
                },
            ]
        },
        {
            id: 3,
            title: 'Зал',
            icon: 'mdi-sofa',
            groups: [
                {
                    id: 1,
                    title: 'group 1',
                    cards: [
                        {
                            domain: 'switch',
                            id: 'switch.hall_sofalamp', 
                            title: 'Диван',
                            state: 'off',
                            icon: 'mdi-lamp'
                        },
                        {
                            domain: 'switch',
                            id: 'switch.hall_outlet2', 
                            title: 'Телевизор',
                            state: 'off',
                            icon: 'mdi-television'
                        },
                        {
                            domain: 'switch',
                            id: 'switch.hall_luster2', 
                            title: 'Большой свет',
                            state: 'off',
                            icon: 'mdi-string-lights'
                        },
                        {
                            domain: 'switch',
                            id: 'switch.hall_luster1', 
                            title: 'Малый свет',
                            state: 'off',
                            icon: 'mdi-vanity-light'
                        },
                    ]
                },
                {
                    id: 2,
                    title: 'group2',
                    cards: [
                        {
                            domain: 'switch',
                            id: 'switch.hall_nightlight', 
                            title: 'Ночник',
                            state: 'off',
                            icon: 'mdi-floor-lamp'
                        },
                        {
                            domain: 'light',
                            id: 'light.hall_pibl',
                            title: 'e-ton',
                            state: 'off',
                            min: 300,
                            max: 900,
                            step: 100,
                            brightness: 0,
                            icon: 'mdi-tablet-dashboard'
                        },
                        {
                            domain: 'spacer',
                            width: 1
                        },
                        {
                            domain: 'gauge',
                            id: 'sensor.hall_temp',
                            title: 'Темп',
                            state: 12.3,
                            unit_of_measurement: '℃',
                            min: 0,
                            max: 40
                        },
                    ]
                }
            ]
        },
        {
            id: 4,
            title: 'Коридор',
            icon: 'mdi-door-open',
            groups: [
                {
                    id: 1,
                    title: 'group 1',
                    cards: [
                        {
                            domain: 'alarm',
                            id: 'alarm.sbu',
                            title: 'СБУ',
                            state: 'disarmed',
                            icon: 'mdi-bell-off-outline'
                        },
                        {
                            domain: 'spacer',
                            width: 0.5
                        },
                        {
                            domain: 'switch',
                            id: 'switch.corridor_light', 
                            title: 'Прихожая',
                            state: 'off',
                            icon: 'mdi-coach-lamp'
                        },
                        {
                            domain: 'light',
                            id: 'light.corridor_wardrobe',
                            title: 'Шкаф',
                            state: 'off',
                            min: 0,
                            max: 100,
                            step: 10,
                            brightness: 0,
                            icon: 'mdi-lava-lamp'
                        },
                        {
                            domain: 'switch',
                            id: 'switch.corridor_passlight', 
                            title: 'Проход',
                            state: 'off',
                            icon: 'mdi-outdoor-lamp'
                        },
                    ]
                }
            ]
        },
        {
            id: 5,
            title: 'Ванная',
            icon: 'mdi-toilet',
            groups: [
                {
                    id: 1,
                    title: 'group 1',
                    cards: [
                        {
                            domain: 'switch',
                            id: 'switch.bathroom_light', 
                            title: 'Свет',
                            state: 'off',
                            icon: 'mdi-toilet'
                        },
                        {
                            domain: 'switch',
                            id: 'switch.bathroom_fan', 
                            title: 'Вентиляция',
                            state: 'off',
                            icon: 'mdi-fan'
                        },
                        {
                            domain: 'switch',
                            id: 'switch.toilet_subscription', 
                            title: 'Подписка',
                            state: 'off',
                            icon: 'mdi-cellphone-message'
                        },
                    ]
                }
            ]
        },
        {
            id: 6,
            title: 'Разное',
            icon: 'mdi-tools',
            groups: [
                {
                    id: 1,
                    title: 'group 1',
                    cards: [
                        {
                            domain: 'binary_sensor',
                            id: 'binary_sensor.pashaxxs',
                            title: 'Паша',
                            icon: 'mdi-account-tie',
                            state: 'off'
                        },
                        {
                            domain: 'binary_sensor',
                            id: 'binary_sensor.xiaomia2',
                            title: 'Таня',
                            icon: 'mdi-account-circle',
                            state: 'off'
                        },
                        {
                            domain: 'binary_sensor',
                            id: 'binary_sensor.krysa',
                            title: 'Аня',
                            icon: 'mdi-account-child-outline',
                            state: 'off'
                        },
                        {
                            domain: 'service',
                            title: 'Пинг',
                            service: 'telegram',
                            method: 'notify',
                            id: 'telegram.ping_me',
                            icon: 'mdi-telegram',
                            state: 'off'
                        },
                        {
                            domain: 'sensor',
                            id: 'sensor.corridor_luminosity',
                            title: 'Освещенность',
                            attributes: {
                                unit_of_measurement: ''
                            },
                            state: 0
                        },
                    ]
                }
            ]
        }
    ]
}
