import { ComponentProps, useContext, useState } from "react";
import { Button, Form } from "react-bootstrap";
import { LanguagesContext } from "../context";

// Language    string  `form:"lang" url:"lang"`
// Name        string  `form:"name" url:"name"`
// Key         string  `form:"key" url:"key"`
// MinPrice    float32 `form:"minPrice,default=-1" url:"minPrice"`
// MaxPrice    float32 `form:"maxPrice,default=-1" url:"maxPrice"`
// Expansion   string  `form:"expansion" url:"expansion"`

interface CardQueryProps extends ComponentProps<"div"> {
    onApply: (query: string) => void
}

const CardQuery = (props: CardQueryProps) => {
    const [foilOnly, setFoilOnly] = useState(false);
    const [inStockOnly, setInStockOnly] = useState(false);
    const [language, setLanguage] = useState('');

    const languages = useContext(LanguagesContext) as LanguageData[];

    const onClearFilters = () => {
        setFoilOnly(false);
        setInStockOnly(false);

        props.onApply('');
    };

    const onApply = () => {
        const lang = languages.find(l => l.longName === language);
        const data = {
            'foilOnly': foilOnly.toString(),
            'inStockOnly': inStockOnly.toString(),
            'lang': lang ? lang.id : '',
        }
        props.onApply(new URLSearchParams(data).toString());
    };

    return (
        <>
            <div className="d-flex my-1 align-items-center">
                <Form.Check
                    type='checkbox'
                    label='Foil'
                    className="mx-1"
                    checked={foilOnly}                
                    onChange={(e) => setFoilOnly(e.target.checked)}
                />
                <Form.Check
                    type='checkbox'
                    label='In stock'
                    className="mx-1"

                    checked={inStockOnly}                
                    onChange={(e) => setInStockOnly(e.target.checked)}
                />
                <Form.Select
                    className="my-2 w-auto mx-4"
                    onChange={e => setLanguage(e.target.value)}
                    value={language}
                >
                    <option>All languages</option> 
                    {languages.map(lang => (
                        <option key={lang.id}>{lang.longName}</option>
                    ))}
                </Form.Select>
            </div>

            <div className="d-flex flex-row-reverse">
                <Button
                    className='mx-1'
                    onClick={onApply}
                >Apply</Button>
                <Button
                    className='mx-1'
                    onClick={onClearFilters}
                >Clear</Button>
            </div>
        </>
    );
}

export default CardQuery;