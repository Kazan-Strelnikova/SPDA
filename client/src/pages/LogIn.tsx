import React, { useState, useEffect, useContext } from 'react';
import { Box, Button, TextField, Typography, Container, Paper, Link } from '@mui/material';
import variables from '../variables.module.scss'; 
import { UserContext } from '../contexts/UserContext';
import { postLoginUser } from '../http/post-login-user';
import { useNavigate } from 'react-router-dom';
import { ErrorOutlineRounded } from '@mui/icons-material';

export const LogInPage: React.FC = () => {
    const navigate = useNavigate();
    const userContext = useContext(UserContext);
    if (!userContext) {
        throw new Error("Log must be used within a UserProvider");
    }
    const { user, setUser } = userContext;

    const [error, setError] = useState<string>('');
    const [formData, setFormData] = useState({
        email: '',
        password: ''
    });

    const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const { name, value } = e.target;
        setFormData({
            ...formData,
            [name]: value
        });
    };

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault();
        if (formData.email === '' || formData.password === ''){
            setError('Все поля должны быть заполнены');
            return;
        }
        const EMAIL_REGEXP = /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/;
        if (!EMAIL_REGEXP.test(formData.email)){
            setError('Некорректный формат почты'); 
            return; 
        } 
        setError('');
        (() => {(async function (){
            try {
                const u = await postLoginUser(formData.email, formData.password)
                setUser(u)
                navigate("/")
            } catch (err: any) {
                setError('Войти не получилось, проверьте введенные данные и повторите'); 
                console.log("login err", err)
            }
            
        }())}) ()


    };
    
    useEffect(() => {
        document.body.style.overflow = 'hidden';
        return () => {
            document.body.style.overflow = 'auto';
        };
    }, []);


    return (
        <Box sx={{ 
            minHeight: '70vh', 
            display: 'flex',
            justifyContent: 'center', 
            alignItems: 'center',
        }}>
            <Container maxWidth="xs" sx={{ my: 4 }}>
                <Box sx={{ 
                    display: 'flex', 
                    flexDirection: 'column',
                    alignItems: 'center',
                    mb: 2
                }}>
                </Box>
                
                <Paper 
                    elevation={3} 
                    sx={{ 
                        p: 4, 
                        borderRadius: 2,
                        width: '100%',
                        backgroundColor: 'white',
                        boxShadow: '0px 4px 20px rgba(0, 0, 0, 0.05)'
                    }}
                >
                    <Box component="form" onSubmit={handleSubmit} noValidate>
                        <TextField
                            margin="normal"
                            // required
                            fullWidth
                            id="email"
                            label="Email"
                            name="email"
                            autoComplete="email"
                            autoFocus
                            value={formData.email}
                            onChange={handleChange}
                            sx={{ mb: 2 }}
                            variant="outlined"
                        />

                        <TextField
                            margin="normal"
                            // required
                            fullWidth
                            name="password"
                            label="Пароль"
                            type="password"
                            id="password"
                            autoComplete="current-password"
                            value={formData.password}
                            onChange={handleChange}
                            sx={{ mb: 3 }}
                        />

                        <Button
                            type="submit"
                            fullWidth
                            variant="contained"
                            sx={{ 
                                py: 1.5, 
                                backgroundColor: variables.primary,
                                '&:hover': {
                                    backgroundColor:variables.dark
                                },
                                mb: 2
                            }}
                        >
                            Войти
                        </Button>
                        
                        {error && 
                        <Box sx={{ 
                            display: 'flex', 
                            flexDirection: 'row',
                            alignItems: 'center',
                            gap: 1,
                        }}>
                            <ErrorOutlineRounded  style={{color: variables.error}}/>
                            <Typography variant='body1' style={{color: variables.error}} >{error}</Typography>
                        </Box>}

                        <Box sx={{ textAlign: 'center', mt: 1 }}>
                            <Typography variant="body2">
                                Нет аккаунта?{' '}
                                <Link 
                                    href="/signup" 
                                    variant="body2" 
                                    sx={{ 
                                        color: variables.dark,
                                        textDecoration: 'none',
                                        fontWeight: 500,
                                        '&:hover': {
                                            textDecoration: 'underline', 
                                        }
                                    }}
                                >
                                    Зарегистрироваться
                                </Link>
                            </Typography>
                        </Box>
                    </Box>
                </Paper>
            </Container>
        </Box>
    );
};