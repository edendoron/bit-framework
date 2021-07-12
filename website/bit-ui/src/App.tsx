import React, {useState} from 'react';
import './App.css';
import {getReports, getBitStatus} from './utils/queryAPI';
import {
    Box,
    Button,
    Card,
    CardContent,
    CardHeader,
    createMuiTheme,
    createStyles,
    Grid,
    makeStyles, MuiThemeProvider,
} from '@material-ui/core'
import {Selector} from "./components/Selector";
import {DatePicker} from "./components/DatePicker";
import {ReportTable} from "./components/ReportTable";
import {StatusTable} from "./components/StatusTable";

const queryTypes = ['Reports', 'BIT Status'];
const userGroups = ['group1', 'group2', 'group3', 'group4', 'groupRafael', 'TemperatureCelsius group', 'group general', 'groupField'];
const filterOptions = ['time', 'tag', 'field'];

const useStyles = makeStyles(() =>
    createStyles({
        dateGrid: {
            marginTop: 20,
        },
        paper: {
            width: '50%',
            marginLeft: '25%',
            marginTop: 20,
        },
        sendButton: {
            border: '3px solid #D48166 ',
            marginBottom: '10px',
        },
    },
));


export const App = () => {
    const classes = useStyles();

    const [queryType, setQueryType] = useState('');
    const [userGroup, setUserGroup] = useState('');
    const [filter, setFilter] = useState('');
    const [isDisabled, setDisabled] = useState(true);
    const [startTime, setStartTime] = useState(new Date());
    const [endTime, setEndTime] = useState(new Date());
    const [data, setData] = useState<string>();

    const changeQueryType = (event: React.ChangeEvent<{ value: unknown }>) => {
        setQueryType(event.target.value as string);
        setFilter('');
        setData('');
    }

    const changeUserGroup = (event: React.ChangeEvent<{ value: unknown }>) => {
        setUserGroup(event.target.value as string);
        setFilter('');
        setQueryType('');
        setDisabled(false);
    }

    const changeFilter = (event: React.ChangeEvent<{ value: unknown }>) => {
        setFilter((event.target.value as string));
    }

    const changeStartTime = (date: Date) => {
        setStartTime(date)
    }

    const changeEndTime = (date: Date) => {
        setEndTime(date)
    }

    const renderData = async () => {
        switch (queryType){
            case 'Reports':
                setData(await getReports(filter, startTime, endTime));
                break;
            case 'BIT Status':
                setData(await getBitStatus(userGroup, startTime, endTime, filter));
                break;
        }
    }

    const theme = createMuiTheme({
        overrides: {
            MuiFilledInput: {
                root: {
                    backgroundColor: 'unset'
                }
            }
        }
    })
    // const theme = createMuiTheme({
    //     palette: {
    //         primary: {
    //             main: CURRENT_THEME.main
    //         }
    //     },
    //     overrides: {
    //         MuiPickersToolbar: {
    //             toolbar: {
    //                 backgroundColor: CURRENT_THEME.el3
    //             }
    //         },
    //         MuiPickersCalendarHeader: {
    //             iconButton: {
    //                 backgroundColor: "transparent",
    //                 color: CURRENT_THEME.main
    //             },
    //             dayLabel: {
    //                 color: CURRENT_THEME.textInv //days in calendar
    //             },
    //             transitionContainer: {
    //                 color: CURRENT_THEME.textInv
    //             }
    //         },
    //         MuiPickersBasePicker: {
    //             pickerView: {
    //                 backgroundColor: CURRENT_THEME.background
    //             }
    //             //   backgroundColor: CURRENT_THEME.el3,
    //         },
    //         MuiPickersDay: {
    //             day: {
    //                 color: CURRENT_THEME.textInv //days in calendar
    //             }
    //             //   daySelected: {
    //             //     backgroundColor: CURRENT_THEME.main, //calendar circle
    //             //   },
    //             //   dayDisabled: {
    //             //     color: CURRENT_THEME.main, // current day
    //             //   },
    //             //   current: {
    //             //     color: CURRENT_THEME.main,
    //             //   },
    //         },
    //         // MuiButton:{
    //         //   textPrimary:{
    //         //     color: CURRENT_THEME.main,
    //         //   }
    //         // },
    //
    //         MuiDialogActions: {
    //             root: {
    //                 backgroundColor: CURRENT_THEME.background
    //             }
    //         },
    //         MuiPickersClock: {
    //             clock: {
    //                 backgroundColor: CURRENT_THEME.el2
    //             }
    //             //   pin:{
    //             //     backgroundColor: CURRENT_THEME.main,
    //             //   },
    //             // },
    //             // MuiPickersClockPointer:{
    //             //   pointer:{
    //             //     backgroundColor: CURRENT_THEME.main,
    //             //   },
    //         },
    //         MuiPickersClockNumber: {
    //             clockNumber: {
    //                 color: CURRENT_THEME.textInv
    //             }
    //             //   numberSelected: {
    //             //     backgroundColor: CURRENT_THEME.main, //calendar circle
    //             //   },
    //         }
    //     }
    // })

    return (
    <MuiThemeProvider theme={theme}>
        <Box bgcolor="#373A36" minHeight="100vh" textAlign="center">
            <Box bgcolor="#D48166" minHeight="100vh" marginRight="150px" marginLeft="150px">
                <Card>
                    <CardHeader title="BIT Framework Query System"/>
                    <CardContent>
                        <Grid container justify='space-evenly'>
                            <Selector menuItems={userGroups} currentValue={userGroup} onChange={changeUserGroup} isDisabled={false}/>
                            <Selector menuItems={queryTypes} currentValue={queryType} onChange={changeQueryType} isDisabled={isDisabled}/>
                            <Selector menuItems={filterOptions} currentValue={filter} onChange={changeFilter} isDisabled={isDisabled}/>
                        </Grid>
                        {filter === 'time' &&
                        <Grid className={classes.dateGrid} container justify='space-evenly'>
                            <DatePicker
                                currentDate={startTime}
                                onDateChange={changeStartTime}
                                placeholder='start time'
                            />
                            <DatePicker
                                currentDate={endTime}
                                onDateChange={changeEndTime}
                                placeholder={'end time'}
                            />
                        </Grid>}
                    </CardContent>
                    <Button className={classes.sendButton} onClick={renderData}>
                        Send
                    </Button>
                </Card>
                {!! data && queryType == 'Reports' && <ReportTable data={JSON.parse(data)}/>}
                {!! data && queryType == 'BIT Status' && <StatusTable data={JSON.parse(data)}/>}
            </Box>
        </Box>
    </MuiThemeProvider>
    );
}

export default App;
