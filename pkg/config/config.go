package config

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type Config struct {
	Logger *logrus.Logger
}

type Loader interface {
	LoadConfig(cfg interface{}) error
}

type loader struct {
	cfg Config
}

func NewConfigLoader(cfg Config) Loader {
	l := &loader{
		cfg: cfg,
	}

	return l
}

func (l *loader) LoadConfig(cfg interface{}) error {
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	err := l.LoadFlags(cfg)
	if err != nil {
		l.cfg.Logger.WithError(err).Errorf("could not load config from command line flags")
		return err
	}

	err = l.LoadFiles(".config", "./config/")
	if err != nil {
		l.cfg.Logger.WithError(err).Errorf("could not load config from file")
		return err
	}

	err = l.LoadEnvs(cfg)
	if err != nil {
		l.cfg.Logger.WithError(err).Errorf("could not load config from environment variables")
		return err
	}

	err = viper.Unmarshal(cfg)
	if err != nil {
		l.cfg.Logger.WithError(err).Errorf("could not unmarshal configuration")
		return err
	}

	return nil
}

func (l *loader) LoadFiles(configName string, paths ...string) error {
	viper.SetConfigName(configName)
	for _, path := range paths {
		viper.AddConfigPath(path)
	}

	err := viper.MergeInConfig()
	if err != nil {
		return err
	}

	return nil
}

func (l *loader) LoadEnvs(cfg interface{}) error {
	t := reflect.TypeOf(cfg)
	err := l.bindEnvs(t)

	return err
}

func (l *loader) bindEnvs(t reflect.Type, parts ...string) error {
	kind := t.Kind()
	if kind == reflect.Pointer || kind == reflect.UnsafePointer || kind == reflect.Interface {
		return l.bindEnvs(t.Elem(), parts...)
	}

	// TODO support kinds as below in field kind switch
	if kind != reflect.Struct {
		return fmt.Errorf("provided type must be a struct or pointing to a struct")
	}

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)

		alias, ok := f.Tag.Lookup("mapstructure")
		if !ok {
			alias = strings.ToLower(f.Name)
		}

		fKind := f.Type.Kind()
		switch fKind {
		case reflect.Chan,
			reflect.Func,
			reflect.Invalid:
			return fmt.Errorf("unsupported kind: %s", fKind.String())
		case reflect.Struct,
			reflect.Pointer,
			reflect.UnsafePointer,
			reflect.Interface:
			err := l.bindEnvs(f.Type, append(parts, alias)...)
			if err != nil {
				return err
			}
		default:
			err := viper.BindEnv(strings.Join(append(parts, alias), "."))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (l *loader) LoadFlags(cfg interface{}) error {
	t := reflect.TypeOf(cfg)
	err := l.bindFlags(t)

	pflag.Parse()

	return err
}

func (l *loader) bindFlags(t reflect.Type, parts ...string) error {
	kind := t.Kind()
	if kind == reflect.Pointer || kind == reflect.UnsafePointer || kind == reflect.Interface {
		return l.bindFlags(t.Elem(), parts...)
	}

	// TODO support kinds as below in field kind switch
	if kind != reflect.Struct {
		return fmt.Errorf("provided type must be a struct or pointing to a struct")
	}

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)

		alias, ok := f.Tag.Lookup("mapstructure")
		if !ok {
			alias = strings.ToLower(f.Name)
		}

		usage, ok := f.Tag.Lookup("usage")
		if !ok {
			usage = "no usage available"
		}

		flagKey := strings.Join(append(parts, alias), "-")

		fKind := f.Type.Kind()
		switch fKind {
		case reflect.Struct,
			reflect.Pointer,
			reflect.UnsafePointer,
			reflect.Interface:
			err := l.bindFlags(f.Type, append(parts, alias)...)
			if err != nil {
				return err
			}
		case reflect.Int:
			//TODO default
			pflag.Int(flagKey, 0, usage)
			err := viper.BindPFlag(flagKey, pflag.Lookup(flagKey))
			if err != nil {
				return err
			}
		case reflect.Int64:
			//TODO default
			pflag.Int64(flagKey, 0, usage)
			err := viper.BindPFlag(flagKey, pflag.Lookup(flagKey))
			if err != nil {
				return err
			}
		case reflect.Int32:
			//TODO default
			pflag.Int32(flagKey, 0, usage)
			err := viper.BindPFlag(flagKey, pflag.Lookup(flagKey))
			if err != nil {
				return err
			}
		case reflect.Int16:
			//TODO default
			pflag.Int16(flagKey, 0, usage)
			err := viper.BindPFlag(flagKey, pflag.Lookup(flagKey))
			if err != nil {
				return err
			}
		case reflect.Int8:
			//TODO default
			pflag.Int8(flagKey, 0, usage)
			err := viper.BindPFlag(flagKey, pflag.Lookup(flagKey))
			if err != nil {
				return err
			}
		case reflect.Uint:
			//TODO default
			pflag.Uint(flagKey, 0, usage)
			err := viper.BindPFlag(flagKey, pflag.Lookup(flagKey))
			if err != nil {
				return err
			}
		case reflect.Uint64:
			//TODO default
			pflag.Uint64(flagKey, 0, usage)
			err := viper.BindPFlag(flagKey, pflag.Lookup(flagKey))
			if err != nil {
				return err
			}
		case reflect.Uint32:
			//TODO default
			pflag.Uint32(flagKey, 0, usage)
			err := viper.BindPFlag(flagKey, pflag.Lookup(flagKey))
			if err != nil {
				return err
			}
		case reflect.Uint16:
			//TODO default
			pflag.Uint16(flagKey, 0, usage)
			err := viper.BindPFlag(flagKey, pflag.Lookup(flagKey))
			if err != nil {
				return err
			}
		case reflect.Uint8:
			//TODO default
			pflag.Uint8(flagKey, 0, usage)
			err := viper.BindPFlag(flagKey, pflag.Lookup(flagKey))
			if err != nil {
				return err
			}
		case reflect.Float64:
			//TODO default
			pflag.Float64(flagKey, 0, usage)
			err := viper.BindPFlag(flagKey, pflag.Lookup(flagKey))
			if err != nil {
				return err
			}
		case reflect.Float32:
			//TODO default
			pflag.Float32(flagKey, 0, usage)
			err := viper.BindPFlag(flagKey, pflag.Lookup(flagKey))
			if err != nil {
				return err
			}
		case reflect.String:
			//TODO default
			pflag.String(flagKey, "", usage)
			err := viper.BindPFlag(flagKey, pflag.Lookup(flagKey))
			if err != nil {
				return err
			}
		case reflect.Bool:
			//TODO default
			pflag.Bool(flagKey, false, usage)
			err := viper.BindPFlag(flagKey, pflag.Lookup(flagKey))
			if err != nil {
				return err
			}
		case reflect.Slice:
			fElemKind := f.Type.Elem().Kind()
			switch fElemKind {
			//TODO more kinds
			case reflect.String:
				//TODO default
				pflag.StringSlice(flagKey, []string{}, usage)
				err := viper.BindPFlag(flagKey, pflag.Lookup(flagKey))
				if err != nil {
					return err
				}
			case reflect.Bool:
				//TODO default
				pflag.BoolSlice(flagKey, []bool{}, usage)
				err := viper.BindPFlag(flagKey, pflag.Lookup(flagKey))
				if err != nil {
					return err
				}
			case reflect.Int:
				//TODO default
				pflag.IntSlice(flagKey, []int{}, usage)
				err := viper.BindPFlag(flagKey, pflag.Lookup(flagKey))
				if err != nil {
					return err
				}
			default:
				return fmt.Errorf("unsupported slice kind: %s", fElemKind)
			}
		default:
			return fmt.Errorf("unsupported kind: %s", fKind)
		}
	}

	return nil
}
