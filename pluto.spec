Name:           pluto
Version:        5.19.0
Release:        1%{?dist}
Summary:        Fairwinds Kubernetes Configuration Validator
License:        Apache-2.0
URL:            https://github.com/FairwindsOps/pluto
Source0:        https://github.com/FairwindsOps/pluto/archive/v%{version}.tar.gz

BuildRequires: make, git, golang, wget

# there's no debug files in this build
%define debug_package %{nil}

%description
Pluto is a tool for validating Kubernetes configuration files.

%prep
%autosetup -n %{name}-%{version}

%build
export PATH=$PWD/go/bin:$PATH
go version
make %{?_smp_mflags} build

%install
install -D -m 0755 %{_builddir}/%{name}-%{version}/pluto %{buildroot}/%{_bindir}/%{name}

%files
%license LICENSE
%doc README.md
%{_bindir}/%{name}

%changelog
* Fri Mar 01 2024 Emanuele Ciurleo <emanuele@ciurleo.com> - 5.19.0-1
- Initial package release
