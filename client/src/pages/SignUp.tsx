import React, { useState, useEffect } from 'react';
import { Box, Button, TextField, Typography, Container, Paper, Link } from '@mui/material';
import variables from '../variables.module.scss'; 

export const SignUpPage: React.FC = () => {
    const [formData, setFormData] = useState({
        email: '',
        password: '',
        confirmPassword: ''
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
        console.log('Form submitted:', formData);
        // registration logic 
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