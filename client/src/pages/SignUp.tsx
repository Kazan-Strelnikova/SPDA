import React, { useState, useEffect, useContext } from 'react';
import { Box, Button, TextField, Typography, Container, Paper, Link } from '@mui/material';
import variables from '../variables.module.scss';
import { ErrorOutlineRounded } from '@mui/icons-material';
import { postRegisterUser } from '../http/post-register-user';
import { useNavigate } from 'react-router-dom';
import { UserContext } from '../contexts/UserContext';

export const SignUpPage: React.FC = () => {
    const navigate = useNavigate()
    const userContext = useContext(UserContext);
    if (!userContext) {
        throw new Error("Log must be used within a UserProvider");
    }
    const { user, setUser } = userContext;

    const [formData, setFormData] = useState({
        name: '',
        lastname: '',
        email: '',
        password: '',
        confirmPassword: ''
    });

    const [error, setError] = useState<string>('');

    const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const { name, value } = e.target;
        setFormData({
            ...formData,
            [name]: value
        });
    };

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault();
        if (formData.name === '' || formData.lastname === '' || formData.email === '' || formData.password === '' || formData.confirmPassword === ''){
            setError('Все поля должны быть заполнены');
            return;
        }
        const EMAIL_REGEXP = /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/;
        if (!EMAIL_REGEXP.test(formData.email)){
            setError('Некорректный формат почты'); 
            return; 
        } 
        if (formData.password !== formData.confirmPassword){
            setError('Пароли не совпадают');
        }
        setError('');
        
        (async () => {
            try {
                const createdUser = await postRegisterUser(
                    formData.name,
                    formData.lastname,
                    formData.email,
                    formData.password
                );

                setUser(createdUser)
                navigate("/") 
            } catch (err: any) {   
                setError("Регистрация не удалась, попробуйте снова");
            }
        }) ()
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
            <Container maxWidth="sm" sx={{ my: 4 }}>
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
                        <Box sx={{ 
                            display: 'flex', 
                            flexDirection: 'row',
                            alignItems: 'center',
                            gap: 2,
                        }}>
                            <TextField
                                margin="normal" 
                                fullWidth
                                id="name"
                                label="Имя"
                                name="name"
                                autoComplete="name"
                                autoFocus
                                value={formData.name}
                                onChange={handleChange}
                                sx={{ mb: 2 }}
                                variant="outlined"
                            />
                            <TextField
                                margin="normal" 
                                fullWidth
                                id="lastname"
                                label="Фамилия"
                                name="lastname"
                                autoComplete="lastname"
                                autoFocus
                                value={formData.lastname}
                                onChange={handleChange}
                                sx={{ mb: 2 }}
                                variant="outlined"
                            />
                        </Box>
                    
                        <TextField
                            margin="normal" 
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
                            autoComplete="new-password"
                            value={formData.password}
                            onChange={handleChange}
                            sx={{ mb: 2 }}
                        />

                        <TextField
                            margin="normal"
                            // required
                            fullWidth
                            name="confirmPassword"
                            label="Подтвердите пароль"
                            type="password"
                            id="confirmPassword"
                            autoComplete="new-password"
                            value={formData.confirmPassword}
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
                                    backgroundColor: variables.dark                                },
                                mb: 2
                            }}
                        >
                            Зарегистрироваться
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
                            Уже есть аккаунт?{' '}
                                <Link 
                                    href="/login" 
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
                                    Войти
                                </Link>
                            </Typography>
                        </Box>
                    </Box>
                </Paper>
            </Container>
        </Box>
    );
};