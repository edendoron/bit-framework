import React, {FC} from 'react';
import {MuiPickersUtilsProvider, KeyboardDateTimePicker} from "@material-ui/pickers";
import DateFnsUtils from "@date-io/date-fns";

interface DatePickerProps {
    currentDate: Date,
    onDateChange: (date: Date) => void,
    label: string,
    disabled: boolean,
}

export const DatePicker: FC<DatePickerProps> = ({currentDate, onDateChange, label, disabled}) => {

    return (
        <MuiPickersUtilsProvider utils={DateFnsUtils}>
            <KeyboardDateTimePicker
                variant="inline"
                ampm={false}
                label={label}
                value={disabled ? null : currentDate}
                disabled={disabled}
                onChange={(date) => onDateChange(date as Date)}
                onError={console.log}
                format="yyyy/MM/dd HH:mm:ss"
                autoOk={true}
                disableFuture={true}
                disableToolbar={true}
                clearable
            />
        </MuiPickersUtilsProvider>
    )
}