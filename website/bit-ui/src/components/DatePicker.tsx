import React, {FC} from 'react';
import {createStyles, makeStyles, TextField, Theme} from "@material-ui/core";

interface DatePickerProps {
    currentDate: Date,
    onChange: (event: React.ChangeEvent<{ value: unknown }>) => void,
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

        },
    }),
);

export const DatePicker: FC<DatePickerProps> = ({currentDate, onChange, placeholder}) => {
    const classes = useStyles();

    return (
        <form className={classes.container} noValidate>
            <TextField
                id="datetime-local"
                label={placeholder}
                type="datetime-local"
                defaultValue={new Date()}
                className={classes.textField}
                InputLabelProps={{
                    shrink: true,
                }}
                value={currentDate}
                onChange={onChange}
            />
        </form>
    )
}