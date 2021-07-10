import React, {FC} from 'react';
import {createStyles, makeStyles, TextField, Theme} from "@material-ui/core";
import {MuiPickersUtilsProvider, KeyboardDateTimePicker} from "@material-ui/pickers";
import DateFnsUtils from "@date-io/date-fns";

interface DatePickerProps {
    currentDate: Date,
    onDateChange: (date: Date) => void,
    placeholder: string,
}

const useStyles = makeStyles((theme: Theme) =>
    createStyles({
        container: {
            display: 'flex',
            flexWrap: 'wrap',
        },
        textField: {
            marginLeft: theme.spacing(1),
            marginRight: theme.spacing(1),
            width: 200,
            color: '#D48166',
        },
    }),
);

export const DatePicker: FC<DatePickerProps> = ({currentDate, onDateChange, placeholder}) => {
    const classes = useStyles();

    return (
        <MuiPickersUtilsProvider utils={DateFnsUtils}>
            <KeyboardDateTimePicker
                variant="inline"
                ampm={false}
                label={placeholder}
                value={currentDate}
                onChange={(date) => onDateChange(date as Date)}
                onError={console.log}
                format="yyyy/MM/dd HH:mm:ss"
            />
        </MuiPickersUtilsProvider>
    )
}