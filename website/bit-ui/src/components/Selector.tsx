import React, {FC} from 'react';
import {createStyles, FormControl, makeStyles, MenuItem, Select} from "@material-ui/core";

interface SelectorProps {
    menuItems: string[],
    currentValue: string,
    onChange: (event: React.ChangeEvent<{ value: unknown }>) => void,
    isDisabled: boolean,
}

const useStyles = makeStyles(() =>
    createStyles({
        formControl: {
            width: '30%',
            alignSelf: 'center',
        },
        select: {
            '&:before': {
                borderColor: '#D48166',
            },
            '&:after': {
                borderColor: '#373A36',
            }
        },
    }),
);

export const Selector: FC<SelectorProps> = ({menuItems, currentValue, onChange, isDisabled}) => {
    const classes = useStyles();

    const renderMenuItems = () => {
        return menuItems.map(item => <MenuItem value={item}>{item}</MenuItem>)
    }
    return (
        <FormControl className={classes.formControl} disabled={isDisabled}>
            <Select
                className={classes.select}
                value={currentValue}
                onChange={onChange}
            >
                {renderMenuItems()}
            </Select>
        </FormControl>
        )
}